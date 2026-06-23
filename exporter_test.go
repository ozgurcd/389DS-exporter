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
	v := reflect.ValueOf(d)
	attrs := make([]*ldap.EntryAttribute, len(metricDefs))
	for i, m := range metricDefs {
		attrs[i] = &ldap.EntryAttribute{
			Name:   m.ldapName,
			Values: []string{fmtFloat(v.Field(m.fieldIdx).Float())},
		}
	}
	return attrs
}

func fmtFloat(f float64) string {
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
			return nil, errors.New("should not reach here")
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
					DN:         "cn=monitor",
					Attributes: attrsToLDAP(obj.DSData{}),
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
	defer saveGlobals()()

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
