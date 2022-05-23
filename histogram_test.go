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

func TestTrapMetrics_HistogramRecordValue(t *testing.T) {
	tests := []struct {
		name        string
		metricName  string
		wantJSON    string
		metricTags  Tags
		metricValue float64
		wantErr     bool
	}{
		{
			name:        "valid",
			metricName:  "test",
			metricValue: 3.14,
			metricTags: Tags{
				Tag{
					Category: "foo",
					Value:    "bar",
				},
			},
			wantJSON: `"_value":"AAEfAAAB"`,
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

			if err = tm.HistogramRecordValue(tt.metricName, tt.metricTags, tt.metricValue); (err != nil) != tt.wantErr {
				t.Errorf("TrapMetrics.HistogramRecordValue() error = %v, wantErr %v", err, tt.wantErr)
			}

			if jm, err := tm.JSONMetrics(); err != nil {
				t.Fatalf("flushing metrics: %s", err)
			} else if !strings.Contains(string(jm), tt.wantJSON) {
				t.Errorf("json metrics want [%v] got [%v]", tt.wantJSON, string(jm))
			}
		})
	}
}

func TestTrapMetrics_HistogramRecordDuration(t *testing.T) {
	tests := []struct {
		name        string
		metricName  string
		wantJSON    string
		metricTags  Tags
		metricValue time.Duration
		wantErr     bool
	}{
		{
			name:        "valid",
			metricName:  "test",
			metricValue: 3 * time.Millisecond,
			metricTags: Tags{
				Tag{
					Category: "foo",
					Value:    "bar",
				},
			},
			wantJSON: `"_value":"AAEe/QAB"`,
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

			if err := tm.HistogramRecordDuration(tt.metricName, tt.metricTags, tt.metricValue); (err != nil) != tt.wantErr {
				t.Errorf("TrapMetrics.HistogramRecordDuration() error = %v, wantErr %v", err, tt.wantErr)
			}

			if jm, err := tm.JSONMetrics(); err != nil {
				t.Fatalf("flushing metrics: %s", err)
			} else if !strings.Contains(string(jm), tt.wantJSON) {
				t.Errorf("json metrics want [%v] got [%v]", tt.wantJSON, string(jm))
			}
		})
	}
}

func TestTrapMetrics_HistogramRecordCountForValue(t *testing.T) {
	tests := []struct {
		name        string
		metricName  string
		wantJSON    string
		metricTags  Tags
		metricValue float64
		metricCount int64
		wantErr     bool
	}{
		{
			name:        "valid",
			metricName:  "test",
			metricValue: 6.2,
			metricCount: 7,
			metricTags: Tags{
				Tag{
					Category: "foo",
					Value:    "bar",
				},
			},
			wantJSON: `"_value":"AAE+AAAH"`,
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

			if err := tm.HistogramRecordCountForValue(tt.metricName, tt.metricTags, tt.metricCount, tt.metricValue); (err != nil) != tt.wantErr {
				t.Errorf("TrapMetrics.HistogramRecordCountForValue() error = %v, wantErr %v", err, tt.wantErr)
			}

			if jm, err := tm.JSONMetrics(); err != nil {
				t.Fatalf("flushing metrics: %s", err)
			} else if !strings.Contains(string(jm), tt.wantJSON) {
				t.Errorf("json metrics want [%v] got [%v]", tt.wantJSON, string(jm))
			}
		})
	}
}
