package trapmetrics

import (
	"context"
	"fmt"

	"github.com/circonus-labs/go-apiclient"
)

func (tm *TrapMetrics) QueueCheckTag(key, val string) {
	if key != "" && val != "" {
		tm.checkTags[key] = val
	}
}

func (tm *TrapMetrics) UpdateCheckTags(ctx context.Context) (*apiclient.CheckBundle, error) {
	if len(tm.checkTags) == 0 {
		return nil, nil
	}

	tags := make([]string, 0, len(tm.checkTags))
	for k, v := range tm.checkTags {
		tags = append(tags, k+":"+v)
	}

	b, err := tm.trap.UpdateCheckTags(ctx, tags)

	// reset the tags
	tm.checkTags = make(map[string]string)

	if err != nil {
		return nil, fmt.Errorf("updating check tags: %w", err)
	}

	if b != nil {
		// if we got a bundle back, tags were updated
		// return it so it can be saved in the cache
		return b, nil
	}

	return nil, nil
}
