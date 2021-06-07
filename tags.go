// Copyright (c) 2021 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package trapmetrics

import (
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
)

type Tag struct {
	Category string
	Value    string
}

type Tags []Tag

// String returns a string representation of tag
func (t *Tag) String() string {
	if t.Category == "" {
		return ""
	}

	return normalizeCategory(t.Category) + ":" + t.Value
}

// Encode returns a base64 encoded tag
func (t *Tag) Encode() string {
	if t.Category == "" {
		return ""
	}

	c := normalizeCategory(t.Category)
	v := t.Value

	encodeFmt := `b"%s"`
	encodedSig := `b"` // has cat or val been previously (or manually) base64 encoded and formatted

	if c != "" && !strings.HasPrefix(c, encodedSig) {
		c = fmt.Sprintf(encodeFmt, base64.StdEncoding.EncodeToString([]byte(c)))
	}
	if v != "" && !strings.HasPrefix(v, encodedSig) {
		v = fmt.Sprintf(encodeFmt, base64.StdEncoding.EncodeToString([]byte(v)))
	}

	return c + ":" + v
}

func normalizeCategory(c string) string {
	return strings.ToLower(strings.ReplaceAll(c, " ", "_"))
}

// String returns a sorted, string list representation of tags
func (tt *Tags) String() string {
	sort.SliceStable(*tt, func(i, j int) bool {
		return (*tt)[i].String() < (*tt)[j].String()
	})

	tags := make([]string, 0, len(*tt))
	for _, t := range *tt {
		tag := t.String()
		if tag != "" {
			tags = append(tags, tag)
		}
	}

	return strings.Join(tags, ",")
}

// Stream returns a base64 encoded string list representation of tags
func (tt *Tags) Encode() string {
	sort.SliceStable(*tt, func(i, j int) bool {
		return (*tt)[i].String() < (*tt)[j].String()
	})

	tags := make([]string, 0, len(*tt))
	for _, t := range *tt {
		tag := t.Encode()
		if tag != "" {
			tags = append(tags, tag)
		}
	}

	return strings.Join(tags, ",")
}

// Stream returns a streamtag encoded string list representation of tags
func (tt *Tags) Stream() string {
	t := tt.Encode()
	if t == "" {
		return t
	}

	return "|ST[" + t + "]"
}
