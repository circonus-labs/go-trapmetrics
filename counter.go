// Copyright (c) 2021 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package trapmetrics

import (
	"fmt"
)

// Note: counters don't take timestamps as they are mutable (e.g. CounterIncrement)
//       rather than track a timestamp separately, when they are flushed the
//       current timestamp is used.

// CounterIncrement will increment the named counter by 1.
func (tm *TrapMetrics) CounterIncrement(name string, tags Tags) error {
	return tm.CounterIncrementByValue(name, tags, 1)
}

// CounterIncrementByValue will increment the named counter by the passed value.
func (tm *TrapMetrics) CounterIncrementByValue(name string, tags Tags, val uint64) error {
	mt := mtCounter

	metricID, err := generateMetricID(name, mt, tags)
	if err != nil {
		return err
	}

	tm.metricsmu.Lock()
	defer tm.metricsmu.Unlock()

	if m, ok := tm.metrics[metricID]; ok {
		if m.Mtype != mtCounter {
			return fmt.Errorf("(%s %s) exists with different type (counter) vs (%s)", name, tags.String(), m.Mtype)
		}
		v, ok := m.Samples[0].(int64)
		if ok {
			v += int64(val)
			m.Samples[0] = v
		}
		return nil
	}

	m, err := tm.newMetric(name, mt, tags)
	if err != nil {
		return fmt.Errorf("(%s %s) failed to initialize (counter): %w", name, tags.String(), err)
	}
	m.Rtype = rtInt64
	m.Samples[0] = int64(val)

	tm.metrics[metricID] = m

	return nil
}

// CounterAdjustByValue will adjust the named counter by the passed value.
func (tm *TrapMetrics) CounterAdjustByValue(name string, tags Tags, val int64) error {
	mt := mtCounter

	metricID, err := generateMetricID(name, mt, tags)
	if err != nil {
		return err
	}

	tm.metricsmu.Lock()
	defer tm.metricsmu.Unlock()

	if m, ok := tm.metrics[metricID]; ok {
		if m.Mtype != mtCounter {
			return fmt.Errorf("(%s %s) exists with different type (counter) vs (%s)", name, tags.String(), m.Mtype)
		}
		v, ok := m.Samples[0].(int64)
		if ok {
			v += val
			m.Samples[0] = v
		}
		return nil
	}

	m, err := tm.newMetric(name, mt, tags)
	if err != nil {
		return fmt.Errorf("(%s %s) failed to initialize (counter): %w", name, tags.String(), err)
	}
	m.Rtype = rtInt64
	m.Samples[0] = val

	tm.metrics[metricID] = m

	return nil
}

// CounterFetch will return the metric identified by name and tags.
func (tm *TrapMetrics) CounterFetch(name string, tags Tags) (*Metric, error) {
	metricID, err := generateMetricID(name, mtCounter, tags)
	if err != nil {
		return nil, err
	}

	tm.metricsmu.Lock()
	defer tm.metricsmu.Unlock()

	if m, ok := tm.metrics[metricID]; ok {
		return m, nil
	}

	return nil, fmt.Errorf("counter %d (%s %s) not found", metricID, name, tags.String())
}
