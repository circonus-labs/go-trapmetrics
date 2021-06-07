// Copyright (c) 2021 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package trapmetrics

import (
	"fmt"
	"time"
)

// TextSet sets a sample with a given timestamp for a text to the passed value
func (tm *TrapMetrics) TextSet(name string, tags Tags, val string, ts *time.Time) error {
	mt := mtText

	metricID, err := generateMetricID(name, mt, tags)
	if err != nil {
		return err
	}
	sampleKey := generateSampleKey(ts)

	tm.metricsmu.Lock()
	defer tm.metricsmu.Unlock()

	if m, ok := tm.metrics[metricID]; ok {
		if m.Mtype != mtText {
			return fmt.Errorf("(%s %s) exists with different type (text) vs (%s)", name, tags.String(), m.Mtype)
		}
		m.Samples[sampleKey] = val
		return nil
	}

	m, err := tm.newMetric(name, mt, tags)
	if err != nil {
		return fmt.Errorf("(%s %s) failed to initialize (text): %w", name, tags.String(), err)
	}
	m.Rtype = rtString
	m.Samples[sampleKey] = val

	tm.metrics[m.ID] = m

	return nil
}

// TextFetch will return the metric identified by name and tags
func (tm *TrapMetrics) TextFetch(name string, tags Tags) (*Metric, error) {
	metricID, err := generateMetricID(name, mtText, tags)
	if err != nil {
		return nil, err
	}

	tm.metricsmu.Lock()
	defer tm.metricsmu.Unlock()

	if m, ok := tm.metrics[metricID]; ok {
		return m, nil
	}

	return nil, fmt.Errorf("text %d (%s %s) not found", metricID, name, tags.String())
}
