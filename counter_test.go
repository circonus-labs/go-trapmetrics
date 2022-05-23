// Copyright (c) 2021 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package trapmetrics

import (
	"strings"
	"testing"
)

func TestTrapMetrics_CounterIncrement(t *testing.T) {
	tests := []struct {
		name        string
		metricName  string
		wantJSON    string
		metricTags  Tags
		metricValue int64
		wantErr     bool
	}{
		{
			name:        "valid",
			metricName:  "test",
			metricValue: 1,
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

			if err = tm.CounterIncrement(tt.metricName, tt.metricTags); (err != nil) != tt.wantErr {
				t.Errorf("TrapMetrics.CounterIncrement() error = %v, wantErr %v", err, tt.wantErr)
			}

			m, err := tm.CounterFetch(tt.metricName, tt.metricTags)
			if (err != nil) != tt.wantErr {
				t.Errorf("TrapMetrics.TextFetch() error = %v, wantErr %v", err, tt.wantErr)
			}

			want := tt.metricValue
			if val, ok := m.Samples[0]; ok {
				if val != want {
					t.Errorf("Invalid value want [%v]%T got [%v]%T", want, want, val, val)
				}
			}

			if jm, err := tm.JSONMetrics(); err != nil {
				t.Fatalf("flushing metrics: %s", err)
			} else if !strings.Contains(string(jm), tt.wantJSON) {
				t.Errorf("json metrics want [%v] got [%v]", tt.wantJSON, string(jm))
			}

		})
	}
}
