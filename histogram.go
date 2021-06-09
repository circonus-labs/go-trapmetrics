// Copyright (c) 2021 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package trapmetrics

import (
	"fmt"
	"time"
)

// Note: histograms don't take timestamps as they already contain multiple samples
//       when they are flushed and serialized the current timestamp is used.

// HistogramRecordTiming adds timing value to histogram.
func (tm *TrapMetrics) HistogramRecordTiming(name string, tags Tags, val float64) error {
	return tm.setValue(name, tags, false, val)
}

// HistogramRecordValue adds value to histogram.
func (tm *TrapMetrics) HistogramRecordValue(name string, tags Tags, val float64) error {
	return tm.setValue(name, tags, false, val)
}

// HistogramRecordDuration adds value to histogram
// (duration is normalized to time.Second, but supports nanosecond granularity).
func (tm *TrapMetrics) HistogramRecordDuration(name string, tags Tags, val time.Duration) error {
	return tm.setDuration(name, tags, false, val)
}

// HistogramRecordCountForValue add count n for value to histogram.
func (tm *TrapMetrics) HistogramRecordCountForValue(name string, tags Tags, count int64, val float64) error {
	return tm.setCountForValue(name, tags, false, count, val)
}

// HistogramFetch will return the metric identified by name and tags.
func (tm *TrapMetrics) HistogramFetch(name string, tags Tags) (*Metric, error) {
	metricID, err := generateMetricID(name, mtHistogram, tags)
	if err != nil {
		return nil, err
	}

	tm.metricsmu.Lock()
	defer tm.metricsmu.Unlock()

	if m, ok := tm.metrics[metricID]; ok {
		return m, nil
	}

	return nil, fmt.Errorf("histogram %d (%s %s) not found", metricID, name, tags.String())
}
