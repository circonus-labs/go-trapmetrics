// Copyright (c) 2021 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package trapmetrics

import "fmt"

//
// Cumulative need to be explicit
//

// Note: histograms don't take timestamps as they already contain multiple samples
//       when they are flushed and serialized the current timestamp is used.

// CumulativeHistogramRecordCountForValue add count n for value to histogram.
func (tm *TrapMetrics) CumulativeHistogramRecordCountForValue(name string, tags Tags, count int64, val float64) error {
	return tm.setCountForValue(name, tags, true, count, val)
}

// CumulativeHistogramFetch will return the metric identified by name and tags.
func (tm *TrapMetrics) CumulativeHistogramFetch(name string, tags Tags) (*Metric, error) {
	metricID, err := generateMetricID(name, mtCumulativeHistogram, tags)
	if err != nil {
		return nil, err
	}

	tm.metricsmu.Lock()
	defer tm.metricsmu.Unlock()

	if m, ok := tm.metrics[metricID]; ok {
		return m, nil
	}

	return nil, fmt.Errorf("cumulative histogram %d (%s %s) not found", metricID, name, tags.String())
}

// CumulativeHistogramTiming adds timing value to histogram
// func (tm *TrapMetrics) CumulativeHistogramRecordTiming(name string, tags Tags, val float64) error {
// 	return tm.setValue(name, tags, true, val)
// }

// CumulativeHistogramRecordValue adds value to histogram
// func (tm *TrapMetrics) CumulativeHistogramRecordValue(name string, tags Tags, val float64) error {
// 	return tm.setValue(name, tags, true, val)
// }

// CumulativeHistogramRecordDuration adds value to histogram
// (duration is normalized to time.Second, but supports nanosecond granularity)
// func (tm *TrapMetrics) CumulativeHistogramRecordDuration(name string, tags Tags, val time.Duration) error {
// 	return tm.setDuration(name, tags, true, val)
// }
