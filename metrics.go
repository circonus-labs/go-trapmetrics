// Copyright (c) 2021 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package trapmetrics

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"io"
	"strings"
	"time"

	"github.com/openhistogram/circonusllhist"
)

const (
	// mytpes (generic).
	mtNone                = "none" //nolint:deadcode,varcheck // metric will get dropped
	mtCounter             = "counter"
	mtGauge               = "gauge"
	mtHistogram           = "histogram"
	mtCumulativeHistogram = "cumulative_histogram"
	mtText                = "text"
	// rtypes (reconnoiter).
	rtInt32               = "i"
	rtUint32              = "I"
	rtInt64               = "l"
	rtUint64              = "L"
	rtFloat64             = "n"
	rtHistogram           = "h"
	rtCumulativeHistogram = "H"
	rtString              = "s"

	// NOTE: max tags and metric name len are enforced here so that
	// details on which metric(s) can be logged. Otherwise, any
	// metric(s) exceeding the limits are rejected by the broker
	// without details on which metric(s) caused the error(s) and
	// what the error(s) actually were, in addition, the broker
	// rejects all metrics sent with the offending metric(s).

	// MaxTags reconnoiter will accept in stream tagged metric name.
	maxTags = 256 // sync w/MAX_TAGS https://github.com/circonus-labs/reconnoiter/blob/master/src/noit_metric.h#L41

	// MaxMetricNameLen reconnoiter will accept (name+stream tags).
	maxMetricNameLen = 4096 // sync w/MAX_METRIC_TAGGED_NAME https://github.com/circonus-labs/reconnoiter/blob/master/src/noit_metric.h#L40
)

var quoteReplacer = strings.NewReplacer(
	`“`, `"`, // smart left double
	`”`, `"`, // smart right double
	`‘`, `'`, // smart left single
	`’`, `'`, // smart right single
)

type Samples map[uint64]interface{}

type Metrics map[uint64]*Metric

type Metric struct {
	Samples Samples
	Name    string
	Mtype   string // set by interface methods
	Rtype   string // set by interface methods
	Tags    Tags
	ID      uint64
}

func (m *Metric) String() string {
	return fmt.Sprintf("id: %d, name: %s, mtype: %s, rtype: %s, tags: %s, samples: %v",
		m.ID,
		m.Name,
		m.Mtype,
		m.Rtype,
		m.Tags.String(),
		m.Samples)
}

func (tm *TrapMetrics) newMetric(metricName, metricType string, tags Tags) (*Metric, error) {
	if metricName == "" {
		return nil, fmt.Errorf("invalid metric name (empty)")
	}
	if metricType == "" {
		return nil, fmt.Errorf("invalid metric type (empty)")
	}
	if len(tags) > maxTags {
		return nil, fmt.Errorf("invalid tags (%d > %d)", len(tags), maxTags)
	}

	id, err := generateMetricID(metricName, metricType, tags)
	if err != nil {
		return nil, err
	}

	m := &Metric{
		ID:      id,
		Name:    metricName,
		Tags:    tags,
		Mtype:   metricType,
		Samples: make(Samples),
	}

	return m, nil
}

func generateMetricID(metricName, metricType string, tags Tags) (uint64, error) {
	h := fnv.New64a()
	_, err := h.Write([]byte(fmt.Sprintf("%s|%s|%s", metricName, metricType, tags.String())))
	if err != nil {
		return 0, fmt.Errorf("hashing name: %w", err)
	}
	return h.Sum64(), nil
}

// generateSampleKey returns a time as a timestamp
// in milliseconds the broker can digest, uses
// current time if ts is nil.
func generateSampleKey(ts *time.Time) uint64 {
	if ts == nil {
		return 0
	}
	return uint64(ts.UTC().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond)))
}

func (tm *TrapMetrics) writeJSONMetrics(w io.Writer) error {
	tm.metricsmu.Lock()
	if len(tm.metrics) == 0 {
		tm.metricsmu.Unlock()
		return nil
	}

	if _, err := w.Write([]byte("{")); err != nil {
		return fmt.Errorf("write {: %w", err)
	}

	var hb bytes.Buffer

	flushTime := time.Now()
	first := true
	for _, m := range tm.metrics {
		tags := m.Tags
		if len(tm.globalTags) > 0 {
			tags = append(tags, tm.globalTags...)
		}
		metricName := m.Name + tags.Stream()
		if len(metricName) > maxMetricNameLen {
			tm.Log.Warnf("metric name exceeds max len (%s)", metricName)
			continue
		}
		brokerType := m.Rtype
		if brokerType == "" {
			tm.Log.Warnf("unknown broker metric type: %s -> %#v", metricName, *m)
			continue
		}

		switch m.Mtype {
		case mtGauge, mtText:
			for sampleKey, sampleValue := range m.Samples {
				_ = writeMetric(w, &first, metricName, brokerType, sampleValue, sampleKey)
			}
		case mtCounter, mtCumulativeHistogram, mtHistogram:
			sampleKey := generateSampleKey(&flushTime)
			if m.Mtype == mtCounter {
				_ = writeMetric(w, &first, metricName, brokerType, m.Samples[0], sampleKey)
			} else {
				hb.Reset()
				if s, ok := m.Samples[0].(*circonusllhist.Histogram); ok {
					if err := s.SerializeB64(&hb); err != nil {
						tm.Log.Warnf("serializing histogram (%s %s): %s", m.Name, m.Tags, err)
						continue
					}
				}
				_ = writeMetric(w, &first, metricName, brokerType, hb.String(), sampleKey)
			}
		}
	}

	if _, err := w.Write([]byte("}")); err != nil {
		return fmt.Errorf("write }: %w", err)
	}

	tm.metrics = make(Metrics)
	tm.metricsmu.Unlock()

	return nil
}

func writeMetric(w io.Writer, first *bool, metricName, metricType string, val interface{}, ts uint64) error {
	value := val

	switch metricType {
	case rtString:
		if s, ok := val.(string); ok {
			// NOTE: convert any 'smart' quotes, escape any embedded quotes, and add string quotes
			value = fmt.Sprintf("%q", quoteReplacer.Replace(s))
		}
	case rtHistogram, rtCumulativeHistogram:
		// NOTE: need to add the string quotes
		value = fmt.Sprintf("%q", val)
	case rtUint64, rtInt64:
		value = fmt.Sprintf(`"%d"`, val)
	case rtFloat64:
		value = fmt.Sprintf(`"%f"`, val)
	}

	// fastest way to get through this
	// otherwise manipulating after
	// to truncate the last comma
	// can take several milliseconds...
	comma := ","
	if *first {
		comma = ""
		*first = false
	}

	metric := fmt.Sprintf(
		`%s%q:{"_type":"%s","_ts":%d,"_value":%v}`,
		comma,
		metricName,
		metricType,
		ts,
		value)

	_, err := io.WriteString(w, metric)

	if err != nil {
		return fmt.Errorf("buf write string: %w", err)
	}

	return nil
}

func (tm *TrapMetrics) jsonMetrics() (bytes.Buffer, error) {
	var buf bytes.Buffer
	buf.Grow(int(tm.bufferSize))
	if err := tm.writeJSONMetrics(&buf); err != nil {
		buf.Reset()
		return buf, fmt.Errorf("writing metrics: %w", err)
	}

	if buf.Len() <= 1 {
		buf.Reset()
		return buf, fmt.Errorf("no valid metrics found")
	}

	return buf, nil
}
