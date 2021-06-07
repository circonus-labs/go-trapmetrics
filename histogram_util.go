// Copyright (c) 2021 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package trapmetrics

import (
	"fmt"
	"time"

	"github.com/openhistogram/circonusllhist"
)

//
// internal histogram support functions - used by both regular and cumulative histograms
//

func (tm *TrapMetrics) setValue(name string, tags Tags, cumulative bool, val float64) error {
	tm.metricsmu.Lock()
	defer tm.metricsmu.Unlock()

	m, err := tm.newHistogram(name, tags, cumulative)
	if err != nil {
		return err
	}

	_ = m.Samples[0].(*circonusllhist.Histogram).RecordValue(val)

	return nil
}

func (tm *TrapMetrics) setDuration(name string, tags Tags, cumulative bool, val time.Duration) error {
	tm.metricsmu.Lock()
	defer tm.metricsmu.Unlock()

	m, err := tm.newHistogram(name, tags, cumulative)
	if err != nil {
		return err
	}

	_ = m.Samples[0].(*circonusllhist.Histogram).RecordDuration(val)

	return nil
}

func (tm *TrapMetrics) setCountForValue(name string, tags Tags, cumulative bool, count int64, val float64) error {
	tm.metricsmu.Lock()
	defer tm.metricsmu.Unlock()

	m, err := tm.newHistogram(name, tags, cumulative)
	if err != nil {
		return err
	}

	_ = m.Samples[0].(*circonusllhist.Histogram).RecordValues(val, count)

	return nil
}

func (tm *TrapMetrics) newHistogram(name string, tags Tags, cumulative bool) (*Metric, error) {
	mt := mtHistogram
	rt := rtHistogram
	if cumulative {
		mt = mtCumulativeHistogram
		rt = rtCumulativeHistogram
	}
	metricID, err := generateMetricID(name, mt, tags)
	if err != nil {
		return nil, err
	}

	if m, ok := tm.metrics[metricID]; ok {
		return m, nil
	}

	m, err := tm.newMetric(name, mt, tags)
	if err != nil {
		return nil, fmt.Errorf("(%s %s) failed to initialize (histogram): %w", name, tags.String(), err)
	}
	m.Rtype = rt
	m.Samples[0] = circonusllhist.New()

	tm.metrics[metricID] = m

	return m, nil
}
