package main

import (
	"reflect"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

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
