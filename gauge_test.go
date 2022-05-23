// Copyright (c) 2021 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package trapmetrics

import (
	"strings"
	"testing"
	"time"
)

func TestTrapMetrics_GaugeSet(t *testing.T) {
	tests := []struct {
		metricValue interface{}
		name        string
		metricName  string
		wantJSON    string
		metricTags  Tags
		wantErr     bool
	}{
		{
			name:        "valid",
			metricName:  "test",
			metricValue: int64(1),
			metricTags: Tags{
				Tag{
					Category: "foo",
					Value:    "bar",
				},
			},
			wantJSON: `"_value":1`,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			tm, err := New(&Config{Trap: FakeTrap{}})
			if err != nil {
				t.Fatalf("unable to initialize TrapMetrics for test: %s", err)
			}

			ts := time.Now()

			if err = tm.GaugeSet(tt.metricName, tt.metricTags, tt.metricValue, &ts); (err != nil) != tt.wantErr {
				t.Errorf("TrapMetrics.GaugeSet() error = %v, wantErr %v", err, tt.wantErr)
			}

			m, err := tm.GaugeFetch(tt.metricName, tt.metricTags)
			if (err != nil) != tt.wantErr {
				t.Errorf("TrapMetrics.TextFetch() error = %v, wantErr %v", err, tt.wantErr)
			}

			want := tt.metricValue
			sk := generateSampleKey(&ts)
			if val, ok := m.Samples[sk]; ok {
				if val != want {
					t.Errorf("Invalid value want [%v]%T got [%v]%T", want, want, val, val)
				}
			} else {
				t.Fatalf("unable to get sample for key %v", sk)
			}

			if jm, err := tm.JSONMetrics(); err != nil {
				t.Fatalf("flushing metrics: %s", err)
			} else if !strings.Contains(string(jm), tt.wantJSON) {
				t.Errorf("json metrics want [%v] got [%v]", tt.wantJSON, string(jm))
			}
		})
	}
}

func TestTrapMetrics_GaugeAdd(t *testing.T) {
	tests := []struct {
		metricValue interface{}
		wantValue   interface{}
		name        string
		metricName  string
		wantJSON    string
		metricTags  Tags
		wantErr     bool
	}{
		{
			name:        "valid int32",
			metricName:  "test",
			metricValue: int32(1),
			metricTags: Tags{
				Tag{
					Category: "foo",
					Value:    "bar",
				},
			},
			wantJSON:  `"_value":2`,
			wantValue: int32(2),
			wantErr:   false,
		},
		{
			name:        "valid int64",
			metricName:  "test",
			metricValue: int64(1),
			metricTags: Tags{
				Tag{
					Category: "foo",
					Value:    "bar",
				},
			},
			wantJSON:  `"_value":2`,
			wantValue: int64(2),
			wantErr:   false,
		},
		{
			name:        "valid uint32",
			metricName:  "test",
			metricValue: uint32(1),
			metricTags: Tags{
				Tag{
					Category: "foo",
					Value:    "bar",
				},
			},
			wantJSON:  `"_value":2`,
			wantValue: uint32(2),
			wantErr:   false,
		},
		{
			name:        "valid uint64",
			metricName:  "test",
			metricValue: uint64(1),
			metricTags: Tags{
				Tag{
					Category: "foo",
					Value:    "bar",
				},
			},
			wantJSON:  `"_value":2`,
			wantValue: uint64(2),
			wantErr:   false,
		},
		{
			name:        "valid float64",
			metricName:  "test",
			metricValue: float64(1.2),
			metricTags: Tags{
				Tag{
					Category: "foo",
					Value:    "bar",
				},
			},
			wantJSON:  `"_value":2.4`,
			wantValue: float64(2.4),
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			tm, err := New(&Config{Trap: FakeTrap{}})
			if err != nil {
				t.Fatalf("unable to initialize TrapMetrics for test: %s", err)
			}

			ts := time.Now()

			if err = tm.GaugeAdd(tt.metricName, tt.metricTags, tt.metricValue, &ts); (err != nil) != tt.wantErr {
				t.Errorf("TrapMetrics.GaugeSet() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err = tm.GaugeAdd(tt.metricName, tt.metricTags, tt.metricValue, &ts); (err != nil) != tt.wantErr {
				t.Errorf("TrapMetrics.GaugeAdd() error = %v, wantErr %v", err, tt.wantErr)
			}

			m, err := tm.GaugeFetch(tt.metricName, tt.metricTags)
			if (err != nil) != tt.wantErr {
				t.Errorf("TrapMetrics.TextFetch() error = %v, wantErr %v", err, tt.wantErr)
			}

			want := tt.wantValue
			sk := generateSampleKey(&ts)
			if val, ok := m.Samples[sk]; ok {
				if val != want {
					t.Errorf("Invalid value want [%v]%T got [%v]%T", want, want, val, val)
				}
			} else {
				t.Fatalf("unable to get sample for key %v", sk)
			}

			if jm, err := tm.JSONMetrics(); err != nil {
				t.Fatalf("flushing metrics: %s", err)
			} else if !strings.Contains(string(jm), tt.wantJSON) {
				t.Errorf("json metrics want [%v] got [%v]", tt.wantJSON, string(jm))
			}
		})
	}
}
