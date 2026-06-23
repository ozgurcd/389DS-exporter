package main

import (
	"testing"

	"github.com/go-ldap/ldap/v3"
	"github.com/ozgurcd/389DS-exporter/obj"
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
	attrs := []*ldap.EntryAttribute{
		{Name: "threads", Values: []string{"1"}},
		{Name: "readwaiters", Values: []string{"2"}},
		{Name: "opsinitiated", Values: []string{"3"}},
		{Name: "opscompleted", Values: []string{"4"}},
		{Name: "dtablesize", Values: []string{"5"}},
		{Name: "anonymousbinds", Values: []string{"6"}},
		{Name: "unauthbinds", Values: []string{"7"}},
		{Name: "simpleauthbinds", Values: []string{"8"}},
		{Name: "strongauthbinds", Values: []string{"9"}},
		{Name: "bindsecurityerrors", Values: []string{"10"}},
		{Name: "inops", Values: []string{"11"}},
		{Name: "readops", Values: []string{"12"}},
		{Name: "compareops", Values: []string{"13"}},
		{Name: "addentryops", Values: []string{"14"}},
		{Name: "removeentryops", Values: []string{"15"}},
		{Name: "modifyentryops", Values: []string{"16"}},
		{Name: "modifyrdnops", Values: []string{"17"}},
		{Name: "searchops", Values: []string{"18"}},
		{Name: "onelevelsearchops", Values: []string{"19"}},
		{Name: "wholesubtreesearchops", Values: []string{"20"}},
		{Name: "referrals", Values: []string{"21"}},
		{Name: "securityerrors", Values: []string{"22"}},
		{Name: "errors", Values: []string{"23"}},
		{Name: "connections", Values: []string{"24"}},
		{Name: "connectionseq", Values: []string{"25"}},
		{Name: "connectionsinmaxthreads", Values: []string{"26"}},
		{Name: "connectionsmaxthreadscount", Values: []string{"27"}},
		{Name: "bytesrecv", Values: []string{"28"}},
		{Name: "bytessent", Values: []string{"29"}},
		{Name: "entriesreturned", Values: []string{"30"}},
		{Name: "referralsreturned", Values: []string{"31"}},
		{Name: "cacheentries", Values: []string{"32"}},
		{Name: "cachehits", Values: []string{"33"}},
	}
	d := parseMonitorAttrs([]*ldap.Entry{{DN: "cn=monitor", Attributes: attrs}})

	want := obj.DSData{
		Threads: 1, Readwaiters: 2, Opsinitiated: 3, Opscompleted: 4,
		Dtablesize: 5, Anonymousbinds: 6, Unauthbinds: 7, Simpleauthbinds: 8,
		Strongauthbinds: 9, Bindsecurityerrors: 10, Inops: 11, Readops: 12,
		Compareops: 13, Addentryops: 14, Removeentryops: 15, Modifyentryops: 16,
		Modifyrdnops: 17, Searchops: 18, Onelevelsearchops: 19, Wholesubtreesearchops: 20,
		Referrals: 21, Securityerrors: 22, Errors: 23, Connections: 24,
		Connectionseq: 25, Connectionsinmaxthreads: 26, Connectionsmaxthreadscount: 27,
		Bytesrecv: 28, Bytessent: 29, Entriesreturned: 30, Referralsreturned: 31,
		Cacheentries: 32, Cachehits: 33,
	}
	if d != want {
		t.Errorf("parseMonitorAttrs mismatch:\ngot  %+v\nwant %+v", d, want)
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
