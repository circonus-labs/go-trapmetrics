// Copyright (c) 2021 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package trapmetrics

import (
	"fmt"
	"time"
)

// GaugeSet sets a sample with a given timestamp for a gauge to the passed value.
func (tm *TrapMetrics) GaugeSet(name string, tags Tags, val interface{}, ts *time.Time) error {
	mt := mtGauge

	metricID, err := generateMetricID(name, mt, tags)
	if err != nil {
		return err
	}

	sampleKey := generateSampleKey(ts)
	rtype := ""

	ok, rt := isValidGaugeType(val)
	if !ok {
		return fmt.Errorf("invalid value for gauge (%v %T)", val, val)
	}
	rtype = rt

	tm.metricsmu.Lock()
	defer tm.metricsmu.Unlock()

	if m, ok := tm.metrics[metricID]; ok {
		if m.Mtype != mtGauge {
			return fmt.Errorf("(%s %s) exists with different type (gauge) vs (%s)", name, tags.String(), m.Mtype)
		}
		m.Samples[sampleKey] = val
		return nil
	}

	m, err := tm.newMetric(name, mt, tags)
	if err != nil {
		return fmt.Errorf("(%s %s) failed to initialize (gauge): %w", name, tags.String(), err)
	}
	m.Rtype = rtype
	m.Samples[sampleKey] = val

	tm.metrics[metricID] = m

	return nil
}

// GaugeAdd adds a sample with a given timestamp for a gauge to the passed value.
func (tm *TrapMetrics) GaugeAdd(name string, tags Tags, val interface{}, ts *time.Time) error {
	mt := mtGauge

	metricID, err := generateMetricID(name, mt, tags)
	if err != nil {
		return err
	}
	sampleKey := generateSampleKey(ts)
	rtype := ""

	ok, rt := isValidGaugeType(val)
	if !ok {
		return fmt.Errorf("invalid value for gauge (%v %T)", val, val)
	}
	rtype = rt

	tm.metricsmu.Lock()
	defer tm.metricsmu.Unlock()

	if m, ok := tm.metrics[metricID]; ok {
		if m.Mtype != mtGauge {
			return fmt.Errorf("(%s %s) exists with different type (gauge) vs (%s)", name, tags.String(), m.Mtype)
		}
		if v, ok := m.Samples[sampleKey]; ok {
			if m.Rtype != rt {
				return fmt.Errorf("(%s %s) exists with different reconnoiter type (%s) vs (%s)", name, tags.String(), m.Rtype, rt)
			}
			m.Samples[sampleKey] = addValByType(m.Rtype, v, val)
		} else {
			m.Samples[sampleKey] = val
		}
		return nil
	}

	m, err := tm.newMetric(name, mt, tags)
	if err != nil {
		return fmt.Errorf("(%s %s) failed to initialize (gauge): %w", name, tags.String(), err)
	}
	m.Rtype = rtype
	m.Samples[sampleKey] = val

	tm.metrics[metricID] = m

	return nil
}

func addValByType(rtype string, base interface{}, val interface{}) interface{} {
	switch rtype {
	case rtInt32:
		return base.(int32) + val.(int32)
	case rtInt64:
		return base.(int64) + val.(int64)
	case rtUint32:
		return base.(uint32) + val.(uint32)
	case rtUint64:
		return base.(uint64) + val.(uint64)
	case rtFloat64:
		return base.(float64) + val.(float64)
	}

	return base
}

func isValidGaugeType(val interface{}) (bool, string) {
	switch val.(type) {
	case int:
		return true, rtInt32
	case int8:
		return true, rtInt32
	case int16:
		return true, rtInt32
	case int32:
		return true, rtInt32
	case int64:
		return true, rtInt64
	case uint:
		return true, rtUint32
	case uint8:
		return true, rtUint32
	case uint16:
		return true, rtUint32
	case uint32:
		return true, rtUint32
	case uint64:
		return true, rtUint64
	case float32:
		return true, rtFloat64
	case float64:
		return true, rtFloat64
	default:
		return false, ""
	}
}

// GaugeFetch will return the metric identified by name and tags.
func (tm *TrapMetrics) GaugeFetch(name string, tags Tags) (*Metric, error) {
	metricID, err := generateMetricID(name, mtGauge, tags)
	if err != nil {
		return nil, err
	}

	tm.metricsmu.Lock()
	defer tm.metricsmu.Unlock()

	if m, ok := tm.metrics[metricID]; ok {
		return m, nil
	}

	return nil, fmt.Errorf("gauge %d (%s %s) not found", metricID, name, tags.String())
}
