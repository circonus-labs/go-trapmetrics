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

func TestTrapMetrics_TextSet(t *testing.T) {
	tests := []struct {
		name        string
		wantJSON    string
		metricName  string
		metricValue string
		metricTags  Tags
		wantErr     bool
	}{
		{
			name:        "valid",
			metricName:  "test",
			metricValue: "test",
			metricTags: Tags{
				Tag{
					Category: "foo",
					Value:    "bar",
				},
			},
			wantJSON: `"_value":"test"`,
			wantErr:  false,
		},
		{
			name:        "valid w/embedded double quotes",
			metricName:  "test",
			metricValue: `test "test"`,
			metricTags: Tags{
				Tag{
					Category: "foo",
					Value:    "bar",
				},
			},
			wantJSON: `"_value":"test \"test\""`,
			wantErr:  false,
		},
		{
			name:        "valid w/embedded single quotes",
			metricName:  "test",
			metricValue: `test 'test'`,
			metricTags: Tags{
				Tag{
					Category: "foo",
					Value:    "bar",
				},
			},
			wantJSON: `"_value":"test 'test'"`,
			wantErr:  false,
		},
		{
			name:        "valid w/embedded 'smart' double quotes",
			metricName:  "test",
			metricValue: `test “smart”`,
			metricTags: Tags{
				Tag{
					Category: "foo",
					Value:    "bar",
				},
			},
			wantJSON: `"_value":"test \"smart\""`,
			wantErr:  false,
		},
		{
			name:        "valid w/embedded 'smart' single quotes",
			metricName:  "test",
			metricValue: `test ‘smart’`,
			metricTags: Tags{
				Tag{
					Category: "foo",
					Value:    "bar",
				},
			},
			wantJSON: `"_value":"test 'smart'"`,
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

			if err = tm.TextSet(tt.metricName, tt.metricTags, tt.metricValue, &ts); (err != nil) != tt.wantErr {
				t.Errorf("TrapMetrics.TextSet() error = %v, wantErr %v", err, tt.wantErr)
			}

			m, err := tm.TextFetch(tt.metricName, tt.metricTags)
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
