package main

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/go-ldap/ldap/v3"
)

func TestParseMonitorAttrs_NoEntries(t *testing.T) {
	d := parseMonitorAttrs(nil)
	if d.Threads != 0 {
		t.Error("expected zero DSData with nil entries")
	}
}

func TestParseMonitorAttrs_Partial(t *testing.T) {
	entries := []*ldap.Entry{{
		DN: "cn=monitor",
		Attributes: []*ldap.EntryAttribute{
			{Name: "threads", Values: []string{"4"}},
			{Name: "readwaiters", Values: []string{"2"}},
		},
	}}
	d := parseMonitorAttrs(entries)
	if d.Threads != 4 {
		t.Errorf("Threads = %v, want 4", d.Threads)
	}
	if d.Readwaiters != 2 {
		t.Errorf("Readwaiters = %v, want 2", d.Readwaiters)
	}
	if d.Cachehits != 0 {
		t.Errorf("Cachehits = %v, want 0 (unset)", d.Cachehits)
	}
}

func TestParseMonitorAttrs_AllFields(t *testing.T) {
	attrs := make([]*ldap.EntryAttribute, len(metricDefs))
	for i, m := range metricDefs {
		attrs[i] = &ldap.EntryAttribute{Name: m.ldapName, Values: []string{strconv.Itoa(i + 1)}}
	}

	d := parseMonitorAttrs([]*ldap.Entry{{DN: "cn=monitor", Attributes: attrs}})

	v := reflect.ValueOf(d)
	for i, m := range metricDefs {
		got := v.Field(m.fieldIdx).Float()
		want := float64(i + 1)
		if got != want {
			t.Errorf("[%d] %s = %v, want %v", m.fieldIdx, m.ldapName, got, want)
		}
	}
}

func TestParseMonitorAttrs_InvalidValuesDefaultToZero(t *testing.T) {
	entries := []*ldap.Entry{{
		DN: "cn=monitor",
		Attributes: []*ldap.EntryAttribute{
			{Name: "threads", Values: []string{"not-a-number"}},
			{Name: "readwaiters", Values: []string{""}},
		},
	}}
	d := parseMonitorAttrs(entries)
	if d.Threads != 0 {
		t.Errorf("Threads = %v, want 0 (invalid parse)", d.Threads)
	}
	if d.Readwaiters != 0 {
		t.Errorf("Readwaiters = %v, want 0 (empty string)", d.Readwaiters)
	}
}

func TestParseMonitorAttrs_MultipleEntries(t *testing.T) {
	entries := []*ldap.Entry{
		{DN: "cn=monitor", Attributes: []*ldap.EntryAttribute{
			{Name: "threads", Values: []string{"8"}},
		}},
		{DN: "cn=other", Attributes: []*ldap.EntryAttribute{
			{Name: "readwaiters", Values: []string{"3"}},
		}},
	}
	d := parseMonitorAttrs(entries)
	if d.Threads != 8 {
		t.Errorf("Threads = %v, want 8", d.Threads)
	}
	if d.Readwaiters != 3 {
		t.Errorf("Readwaiters = %v, want 3", d.Readwaiters)
	}
}

func TestParseFloatWithDefault(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  float64
	}{
		{"valid integer", "42", 42},
		{"valid float", "3.14", 3.14},
		{"valid negative", "-5", -5},
		{"empty string", "", 0},
		{"invalid string", "not-a-number", 0},
		{"scientific notation", "1e5", 100000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseFloatWithDefault(tt.value, tt.name)
			if got != tt.want {
				t.Errorf("parseFloatWithDefault(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
