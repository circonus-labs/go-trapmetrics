// Copyright (c) 2021 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package trapmetrics

import (
	"testing"
)

func TestTag_String(t *testing.T) {
	type fields struct {
		Category string
		Value    string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "valid c:v", fields: fields{Category: "foo", Value: "bar"}, want: "foo:bar"},
		{name: "valid c:", fields: fields{Category: "foo", Value: ""}, want: "foo:"},
		{name: "invalid nocat", fields: fields{Category: "", Value: ""}, want: ""},
		{name: "invalid nocat", fields: fields{Category: "", Value: "bar"}, want: ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tag{
				Category: tt.fields.Category,
				Value:    tt.fields.Value,
			}
			if got := tr.String(); got != tt.want {
				t.Errorf("Tag.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTag_Encode(t *testing.T) {
	type fields struct {
		Category string
		Value    string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "valid c:v", fields: fields{Category: "foo", Value: "bar"}, want: `b"Zm9v":b"YmFy"`},
		{name: "valid c:", fields: fields{Category: "foo", Value: ""}, want: `b"Zm9v":`},
		{name: "invalid nocat", fields: fields{Category: "", Value: ""}, want: ""},
		{name: "invalid nocat", fields: fields{Category: "", Value: "bar"}, want: ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tag{
				Category: tt.fields.Category,
				Value:    tt.fields.Value,
			}
			if got := tr.Encode(); got != tt.want {
				t.Errorf("Tag.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTags_String(t *testing.T) {
	tests := []struct {
		name string
		tt   *Tags
		want string
	}{
		{name: "valid c:v,c:v", tt: &Tags{Tag{Category: "foo", Value: "bar"}, Tag{Category: "baz", Value: "qux"}}, want: `baz:qux,foo:bar`},
		{name: "valid c:,c:v", tt: &Tags{Tag{Category: "foo", Value: "bar"}, Tag{Category: "baz", Value: ""}}, want: `baz:,foo:bar`},
		{name: "valid c:v", tt: &Tags{Tag{Category: "foo", Value: "bar"}, Tag{Category: "", Value: "qux"}}, want: `foo:bar`},
		{name: "invalid nocat", tt: &Tags{Tag{Category: "", Value: ""}, Tag{Category: "", Value: "qux"}}, want: ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tt.String(); got != tt.want {
				t.Errorf("Tags.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTags_Encode(t *testing.T) {
	tests := []struct {
		name string
		tt   *Tags
		want string
	}{
		{name: "valid c:v,c:v", tt: &Tags{Tag{Category: "foo", Value: "bar"}, Tag{Category: "baz", Value: "qux"}}, want: `b"YmF6":b"cXV4",b"Zm9v":b"YmFy"`},
		{name: "valid c:,c:v", tt: &Tags{Tag{Category: "foo", Value: "bar"}, Tag{Category: "baz", Value: ""}}, want: `b"YmF6":,b"Zm9v":b"YmFy"`},
		{name: "valid c:v", tt: &Tags{Tag{Category: "foo", Value: "bar"}, Tag{Category: "", Value: "qux"}}, want: `b"Zm9v":b"YmFy"`},
		{name: "invalid nocat", tt: &Tags{Tag{Category: "", Value: ""}, Tag{Category: "", Value: "qux"}}, want: ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tt.Encode(); got != tt.want {
				t.Errorf("Tags.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTags_Stream(t *testing.T) {
	tests := []struct {
		name string
		tt   *Tags
		want string
	}{
		{name: "valid c:v,c:v", tt: &Tags{Tag{Category: "foo", Value: "bar"}, Tag{Category: "baz", Value: "qux"}}, want: `|ST[b"YmF6":b"cXV4",b"Zm9v":b"YmFy"]`},
		{name: "valid c:,c:v", tt: &Tags{Tag{Category: "foo", Value: "bar"}, Tag{Category: "baz", Value: ""}}, want: `|ST[b"YmF6":,b"Zm9v":b"YmFy"]`},
		{name: "valid c:v", tt: &Tags{Tag{Category: "foo", Value: "bar"}, Tag{Category: "", Value: "qux"}}, want: `|ST[b"Zm9v":b"YmFy"]`},
		{name: "invalid nocat", tt: &Tags{Tag{Category: "", Value: ""}, Tag{Category: "", Value: "qux"}}, want: ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tt.Stream(); got != tt.want {
				t.Errorf("Tags.Stream() = %v, want %v", got, tt.want)
			}
		})
	}
}
