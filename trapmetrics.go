// Copyright (c) 2021 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package trapmetrics

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/circonus-labs/go-trapcheck"
)

const (
	defaultBufferSize = uint(32768)
)

// Trap defines the interface for for submitting metrics
type Trap interface {
	SendMetrics(ctx context.Context, metrics *strings.Builder) (*trapcheck.TrapResult, error)
}

type Config struct {
	// Trap instance of go-trapcheck (or something satisfying Trap interface) to use trapmetrics as a
	// metric container and handle transport externally, pass nil
	Trap Trap

	// Logger instance of something satisfying Logger interface (default: log.Logger with ioutil.Discard)
	Logger Logger

	// GlobalTags is a list of tags to be added to every metric
	GlobalTags Tags

	// BufferSize size of metric buffer (when flushing), default is defaultBufferSize above
	BufferSize uint
}
type TrapMetrics struct {
	trap       Trap
	metrics    Metrics
	Log        Logger
	globalTags Tags
	metricsmu  sync.Mutex
	bufferSize uint
}

func New(cfg *Config) (*TrapMetrics, error) {
	tm := &TrapMetrics{
		trap:       cfg.Trap,
		metrics:    make(Metrics),
		globalTags: cfg.GlobalTags,
	}

	if cfg.Logger != nil {
		tm.Log = cfg.Logger
	} else {
		tm.Log = &LogWrapper{
			Log:   log.New(ioutil.Discard, "", log.LstdFlags),
			Debug: false,
		}
	}
	if cfg.BufferSize == 0 {
		tm.bufferSize = defaultBufferSize
	}

	return tm, nil
}

// JSONMetrics returns the current metrics in JSON format or an error - to be used
// when handling submission of metrics externally (e.g. aggregating multiple sets
// of metrics from different trapmetrics containers)
func (tm *TrapMetrics) JSONMetrics() ([]byte, error) {
	m, err := tm.jsonMetrics()
	if err != nil {
		return []byte{}, err
	}
	return []byte(m.String()), nil
}

type Result struct {
	CheckUUID       string
	Error           string
	SubmitUUID      string
	Filtered        uint64
	Stats           uint64
	SubmitDuration  time.Duration
	LastReqDuration time.Duration
	EncodeDuration  time.Duration
	FlushDuration   time.Duration
	BytesSent       int
}

// Flush sends metrics to the configured trap check, returns result or an error
func (tm *TrapMetrics) Flush(ctx context.Context) (*Result, error) {
	if tm.trap == nil {
		return nil, fmt.Errorf("no trap check configured")
	}

	start := time.Now()
	data, err := tm.jsonMetrics()
	if err != nil {
		return nil, fmt.Errorf("packaging metrics for submission: %w", err)
	}
	if data.Len() == 0 {
		return &Result{Error: "no metrics to send"}, nil
	}
	result := &Result{
		EncodeDuration: time.Since(start),
	}

	smResult, err := tm.trap.SendMetrics(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("submitting metrics to broker: %w", err)
	}

	result.CheckUUID = smResult.CheckUUID
	result.Error = smResult.Error
	result.SubmitUUID = smResult.SubmitUUID
	result.Stats = smResult.Stats
	result.Filtered = smResult.Filtered
	result.SubmitDuration = smResult.SubmitDuration
	result.LastReqDuration = smResult.LastReqDuration
	result.BytesSent = smResult.BytesSent
	result.FlushDuration = time.Since(start)

	tm.Log.Debugf("flush -- C:%s, S:%s, E:%s, Stats:%d, Filtered:%d, Bytes:%d, Encode:%s, Submit:%s, LastReq:%s, Flush:%s",
		result.CheckUUID, result.SubmitUUID, result.Error,
		result.Stats, result.Filtered, result.BytesSent,
		result.EncodeDuration, result.SubmitDuration, result.LastReqDuration, result.FlushDuration)

	return result, nil
}
