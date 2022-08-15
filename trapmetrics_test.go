// Copyright (c) 2021 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package trapmetrics

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/circonus-labs/go-apiclient"
	"github.com/circonus-labs/go-trapcheck"
)

type FakeTrap struct {
}

func (ft FakeTrap) SendMetrics(ctx context.Context, metrics bytes.Buffer) (*trapcheck.TrapResult, error) {
	fmt.Println(metrics)
	return nil, nil
}
func (ft FakeTrap) UpdateCheckTags(ctx context.Context, tags []string) (*apiclient.CheckBundle, error) {
	return nil, nil
}

func TestNew(t *testing.T) {
	tests := []struct {
		cfg     *Config
		name    string
		wantErr bool
	}{
		{
			name:    "invalid",
			cfg:     nil,
			wantErr: true,
		},
		{
			name: "valid",
			cfg: &Config{
				Trap: FakeTrap{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
