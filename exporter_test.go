package main

import (
	"errors"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ozgurcd/389DS-exporter/obj"
)

type mockLDAP struct {
	searchFunc func(*ldap.SearchRequest) (*ldap.SearchResult, error)
	closeFunc  func() error
}

func (m *mockLDAP) Search(req *ldap.SearchRequest) (*ldap.SearchResult, error) {
	return m.searchFunc(req)
}

func (m *mockLDAP) Close() error {
	return m.closeFunc()
}

func saveGlobals() func() {
	origServer := server
	origPort := port
	origTimeout := ldapTimeout
	return func() {
		server = origServer
		port = origPort
		ldapTimeout = origTimeout
	}
}

func attrsToLDAP(d obj.DSData) []*ldap.EntryAttribute {
	return []*ldap.EntryAttribute{
		{Name: "threads", Values: []string{floatStr(d.Threads)}},
		{Name: "readwaiters", Values: []string{floatStr(d.Readwaiters)}},
		{Name: "opsinitiated", Values: []string{floatStr(d.Opsinitiated)}},
		{Name: "opscompleted", Values: []string{floatStr(d.Opscompleted)}},
		{Name: "dtablesize", Values: []string{floatStr(d.Dtablesize)}},
		{Name: "anonymousbinds", Values: []string{floatStr(d.Anonymousbinds)}},
		{Name: "unauthbinds", Values: []string{floatStr(d.Unauthbinds)}},
		{Name: "simpleauthbinds", Values: []string{floatStr(d.Simpleauthbinds)}},
		{Name: "strongauthbinds", Values: []string{floatStr(d.Strongauthbinds)}},
		{Name: "bindsecurityerrors", Values: []string{floatStr(d.Bindsecurityerrors)}},
		{Name: "inops", Values: []string{floatStr(d.Inops)}},
		{Name: "readops", Values: []string{floatStr(d.Readops)}},
		{Name: "compareops", Values: []string{floatStr(d.Compareops)}},
		{Name: "addentryops", Values: []string{floatStr(d.Addentryops)}},
		{Name: "removeentryops", Values: []string{floatStr(d.Removeentryops)}},
		{Name: "modifyentryops", Values: []string{floatStr(d.Modifyentryops)}},
		{Name: "modifyrdnops", Values: []string{floatStr(d.Modifyrdnops)}},
		{Name: "searchops", Values: []string{floatStr(d.Searchops)}},
		{Name: "onelevelsearchops", Values: []string{floatStr(d.Onelevelsearchops)}},
		{Name: "wholesubtreesearchops", Values: []string{floatStr(d.Wholesubtreesearchops)}},
		{Name: "referrals", Values: []string{floatStr(d.Referrals)}},
		{Name: "securityerrors", Values: []string{floatStr(d.Securityerrors)}},
		{Name: "errors", Values: []string{floatStr(d.Errors)}},
		{Name: "connections", Values: []string{floatStr(d.Connections)}},
		{Name: "connectionseq", Values: []string{floatStr(d.Connectionseq)}},
		{Name: "connectionsinmaxthreads", Values: []string{floatStr(d.Connectionsinmaxthreads)}},
		{Name: "connectionsmaxthreadscount", Values: []string{floatStr(d.Connectionsmaxthreadscount)}},
		{Name: "bytesrecv", Values: []string{floatStr(d.Bytesrecv)}},
		{Name: "bytessent", Values: []string{floatStr(d.Bytessent)}},
		{Name: "entriesreturned", Values: []string{floatStr(d.Entriesreturned)}},
		{Name: "referralsreturned", Values: []string{floatStr(d.Referralsreturned)}},
		{Name: "cacheentries", Values: []string{floatStr(d.Cacheentries)}},
		{Name: "cachehits", Values: []string{floatStr(d.Cachehits)}},
	}
}

func floatStr(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func TestNewExporter(t *testing.T) {
	e := NewExporter()
	if e == nil {
		t.Fatal("NewExporter() returned nil")
	}
}

func TestNewExporterDescsNonNil(t *testing.T) {
	e := NewExporter()
	v := reflect.ValueOf(e).Elem()
	descType := reflect.TypeFor[*prometheus.Desc]()
	for i := range v.NumField() {
		f := v.Field(i)
		if f.Type() != descType {
			continue
		}
		if f.IsNil() {
			t.Errorf("%s desc is nil", v.Type().Field(i).Name)
		}
	}
}

func TestDescribeSendsAllDescriptors(t *testing.T) {
	e := NewExporter()
	ch := make(chan *prometheus.Desc, 100)

	e.Describe(ch)
	close(ch)

	var count int
	for range ch {
		count++
	}

	// Expect 33 descriptors (all fields in Exporter struct)
	if count != 33 {
		t.Errorf("Describe sent %d descriptors, want 33", count)
	}
}

func TestGetLDAPConn_DialSuccess(t *testing.T) {
	defer saveGlobals()()
	server = "ldap.example.com"
	port = 389
	ldapTimeout = 5 * time.Second

	e := &Exporter{
		dial: func(addr string) (LDAPClient, error) {
			return &mockLDAP{closeFunc: func() error { return nil }}, nil
		},
	}

	c, err := e.getLDAPConn()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}

	// Second call should return cached conn
	c2, err := e.getLDAPConn()
	if err != nil {
		t.Fatalf("unexpected error on cached call: %v", err)
	}
	if c != c2 {
		t.Error("expected cached connection, got new one")
	}
}

func TestGetLDAPConn_DialError(t *testing.T) {
	defer saveGlobals()()
	server = "ldap.example.com"
	port = 389
	ldapTimeout = 5 * time.Second

	e := &Exporter{
		dial: func(addr string) (LDAPClient, error) {
			return nil, errors.New("dial refused")
		},
	}

	_, err := e.getLDAPConn()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetLDAPConn_NoDial(t *testing.T) {
	e := &Exporter{}
	_, err := e.getLDAPConn()
	if err == nil {
		t.Fatal("expected error when no dial function configured")
	}
}

func TestSearchLDAP_Success(t *testing.T) {
	mock := &mockLDAP{
		searchFunc: func(req *ldap.SearchRequest) (*ldap.SearchResult, error) {
			return &ldap.SearchResult{
				Entries: []*ldap.Entry{{
					DN: "cn=monitor",
					Attributes: []*ldap.EntryAttribute{
						{Name: "threads", Values: []string{"4"}},
					},
				}},
			}, nil
		},
	}

	data, err := searchLDAP(mock, 5*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.Threads != 4 {
		t.Errorf("Threads = %v, want 4", data.Threads)
	}
}

func TestSearchLDAP_Error(t *testing.T) {
	mock := &mockLDAP{
		searchFunc: func(req *ldap.SearchRequest) (*ldap.SearchResult, error) {
			return nil, errors.New("search failed")
		},
	}

	_, err := searchLDAP(mock, 5*time.Second)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSearchLDAP_Timeout(t *testing.T) {
	block := make(chan struct{})
	mock := &mockLDAP{
		searchFunc: func(req *ldap.SearchRequest) (*ldap.SearchResult, error) {
			<-block
			return nil, nil
		},
	}
	defer close(block)

	_, err := searchLDAP(mock, 0)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestCloseLDAPConn(t *testing.T) {
	closed := false
	mock := &mockLDAP{
		searchFunc: func(req *ldap.SearchRequest) (*ldap.SearchResult, error) {
			return nil, nil
		},
		closeFunc: func() error {
			closed = true
			return nil
		},
	}

	e := &Exporter{ldapConn: mock}
	e.closeLDAPConn()

	if !closed {
		t.Error("expected Close() to be called")
	}
	if e.ldapConn != nil {
		t.Error("expected ldapConn to be nil after close")
	}
}

func TestCollect_Success(t *testing.T) {
	defer saveGlobals()()
	server = "ldap.example.com"
	port = 389
	ldapTimeout = 5 * time.Second

	mock := &mockLDAP{
		searchFunc: func(req *ldap.SearchRequest) (*ldap.SearchResult, error) {
			return &ldap.SearchResult{
				Entries: []*ldap.Entry{{
					DN: "cn=monitor",
					Attributes: attrsToLDAP(obj.DSData{
						Threads: 1, Readwaiters: 2, Opsinitiated: 3, Opscompleted: 4,
						Dtablesize: 5, Anonymousbinds: 6, Unauthbinds: 7, Simpleauthbinds: 8,
						Strongauthbinds: 9, Bindsecurityerrors: 10, Inops: 11, Readops: 12,
						Compareops: 13, Addentryops: 14, Removeentryops: 15, Modifyentryops: 16,
						Modifyrdnops: 17, Searchops: 18, Onelevelsearchops: 19, Wholesubtreesearchops: 20,
						Referrals: 21, Securityerrors: 22, Errors: 23, Connections: 24,
						Connectionseq: 25, Connectionsinmaxthreads: 26, Connectionsmaxthreadscount: 27,
						Bytesrecv: 28, Bytessent: 29, Entriesreturned: 30, Referralsreturned: 31,
						Cacheentries: 32, Cachehits: 33,
					}),
				}},
			}, nil
		},
		closeFunc: func() error { return nil },
	}

	e := NewExporter()
	e.dial = func(addr string) (LDAPClient, error) { return mock, nil }

	ch := make(chan prometheus.Metric, 40)
	e.Collect(ch)
	close(ch)

	count := 0
	for range ch {
		count++
	}
	if count != 33 {
		t.Errorf("expected 33 metrics, got %d", count)
	}
}

func TestCollect_SearchError(t *testing.T) {
	defer saveGlobals()()
	server = "ldap.example.com"
	port = 389
	ldapTimeout = 5 * time.Second

	mock := &mockLDAP{
		searchFunc: func(req *ldap.SearchRequest) (*ldap.SearchResult, error) {
			return nil, errors.New("search exploded")
		},
		closeFunc: func() error { return nil },
	}

	e := NewExporter()
	e.dial = func(addr string) (LDAPClient, error) { return mock, nil }

	ch := make(chan prometheus.Metric, 40)
	e.Collect(ch)
	close(ch)

	for range ch {
		t.Error("expected no metrics on search error")
	}
}

func TestCollectHandlesConnectionError(t *testing.T) {
	// Save and restore package-level vars
	origServer := server
	origPort := port
	origTimeout := ldapTimeout
	defer func() {
		server = origServer
		port = origPort
		ldapTimeout = origTimeout
	}()

	server = "192.0.2.1" // TEST-NET address, guaranteed unroutable
	port = 1
	ldapTimeout = 0

	e := NewExporter()
	ch := make(chan prometheus.Metric, 100)

	// Should not panic, should not send any metrics
	e.Collect(ch)
	close(ch)

	// No metrics should be produced when LDAP is unreachable
	count := 0
	for range ch {
		count++
	}
	if count != 0 {
		t.Errorf("Collect produced %d metrics on connection failure, want 0", count)
	}
}
