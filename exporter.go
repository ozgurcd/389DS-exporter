package main

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/prometheus/client_golang/prometheus"
)

type metricKind int

const (
	gaugeKind metricKind = iota
	counterKind
)

type metricDef struct {
	ldapName string
	fieldIdx int
	help     string
	kind     metricKind
	label    string
}

// metricDefs is the single source of truth for all 33 metrics.
// fieldIdx must match obj.DSData struct field order.
var metricDefs = []metricDef{
	{ldapName: "threads", fieldIdx: 0, help: "Number of Threads max configured", kind: gaugeKind, label: "threads"},
	{ldapName: "readwaiters", fieldIdx: 1, help: "Current number of threads waiting to read data from a client", kind: gaugeKind, label: "readwaiters"},
	{ldapName: "opsinitiated", fieldIdx: 2, help: "Current number of operations the server has initiated since it started", kind: counterKind, label: "opsinitiated"},
	{ldapName: "opscompleted", fieldIdx: 3, help: "Current number of operations the server has completed since it started", kind: counterKind, label: "opscompleted"},
	{ldapName: "dtablesize", fieldIdx: 4, help: "The number of file descriptors available to the directory. Essentially, this value shows how many additional concurrent connections can be serviced by the directory", kind: gaugeKind, label: "dtablesize"},
	{ldapName: "anonymousbinds", fieldIdx: 5, help: "Number of Anonymous Binds", kind: counterKind, label: "anonymousbinds"},
	{ldapName: "unauthbinds", fieldIdx: 6, help: "Number of Unauth Binds", kind: counterKind, label: "unauthbinds"},
	{ldapName: "simpleauthbinds", fieldIdx: 7, help: "Number of Simple Auth Binds", kind: counterKind, label: "simpleauthbinds"},
	{ldapName: "strongauthbinds", fieldIdx: 8, help: "Number of Strong Auth Binds", kind: counterKind, label: "strongauthbinds"},
	{ldapName: "bindsecurityerrors", fieldIdx: 9, help: "Number of Bind Security Errors", kind: counterKind, label: "bindsecurityerrors"},
	{ldapName: "inops", fieldIdx: 10, help: "Number of All Requests", kind: counterKind, label: "inops"},
	{ldapName: "readops", fieldIdx: 11, help: "Number of Read Operations", kind: counterKind, label: "readops"},
	{ldapName: "compareops", fieldIdx: 12, help: "Number of Compare Operations", kind: counterKind, label: "compareops"},
	{ldapName: "addentryops", fieldIdx: 13, help: "Number of Add Entry Operations", kind: counterKind, label: "addentryops"},
	{ldapName: "removeentryops", fieldIdx: 14, help: "Number of Remove Entry Operations", kind: counterKind, label: "removeentryops"},
	{ldapName: "modifyentryops", fieldIdx: 15, help: "Number of Modify Entry Operations", kind: counterKind, label: "modifyentryops"},
	{ldapName: "modifyrdnops", fieldIdx: 16, help: "Number of Modify RDN Operations", kind: counterKind, label: "modifyrdnops"},
	{ldapName: "searchops", fieldIdx: 17, help: "Number of LDAP Search Requests", kind: counterKind, label: "searchops"},
	{ldapName: "onelevelsearchops", fieldIdx: 18, help: "Number of one-level Search Requests", kind: counterKind, label: "onelevelsearchops"},
	{ldapName: "wholesubtreesearchops", fieldIdx: 19, help: "Number of subtree-level Search Requests", kind: counterKind, label: "wholesubtreesearchops"},
	{ldapName: "referrals", fieldIdx: 20, help: "Number of LDAP referrals", kind: counterKind, label: "referrals"},
	{ldapName: "securityerrors", fieldIdx: 21, help: "Number of Security Errors", kind: counterKind, label: "securityerrors"},
	{ldapName: "errors", fieldIdx: 22, help: "Number of Errors", kind: counterKind, label: "errors"},
	{ldapName: "connections", fieldIdx: 23, help: "Number of Connections in Open State at the sampling time", kind: gaugeKind, label: "connections"},
	{ldapName: "connectionseq", fieldIdx: 24, help: "Total Number of Connections opened", kind: counterKind, label: "connectionseq"},
	{ldapName: "connectionsinmaxthreads", fieldIdx: 25, help: "Number of connections that are currently in a max thread state", kind: gaugeKind, label: "connectionsinmaxthreads"},
	{ldapName: "connectionsmaxthreadscount", fieldIdx: 26, help: "Number of connectionsmaxthreadscount", kind: gaugeKind, label: "connectionsmaxthreadscount"},
	{ldapName: "bytesrecv", fieldIdx: 27, help: "Total number of bytes received", kind: counterKind, label: "bytesrecv"},
	{ldapName: "bytessent", fieldIdx: 28, help: "Total number of bytes sent", kind: counterKind, label: "bytessent"},
	{ldapName: "entriesreturned", fieldIdx: 29, help: "Number of Entries Returned", kind: counterKind, label: "entriesreturned"},
	{ldapName: "referralsreturned", fieldIdx: 30, help: "Number of Referrals Returned", kind: counterKind, label: "referralsreturned"},
	{ldapName: "cacheentries", fieldIdx: 31, help: "Number of Cache Entries", kind: gaugeKind, label: "cacheentries"},
	{ldapName: "cachehits", fieldIdx: 32, help: "Number of Cache Hits", kind: counterKind, label: "cachehits"},
}

// ldapFieldMap maps LDAP attribute names to DSData field indices.
var ldapFieldMap map[string]int

func init() {
	ldapFieldMap = make(map[string]int, len(metricDefs))
	for _, m := range metricDefs {
		ldapFieldMap[m.ldapName] = m.fieldIdx
	}
}

type LDAPClient interface {
	Search(searchRequest *ldap.SearchRequest) (*ldap.SearchResult, error)
	Close() error
}

type ldapClient struct {
	conn *ldap.Conn
}

func (c *ldapClient) Search(req *ldap.SearchRequest) (*ldap.SearchResult, error) {
	return c.conn.Search(req)
}

func (c *ldapClient) Close() error {
	return c.conn.Close()
}

type dialFunc func(string) (LDAPClient, error)

func runWithTimeout[T any](timeout time.Duration, fn func() (T, error)) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	type result struct {
		val T
		err error
	}
	resultCh := make(chan result, 1)
	go func() {
		v, e := fn()
		resultCh <- result{val: v, err: e}
	}()

	select {
	case r := <-resultCh:
		return r.val, r.err
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()
	}
}

// Exporter stores metrics from 389DS
type Exporter struct {
	mu       sync.Mutex
	ldapConn LDAPClient
	dial     dialFunc
	descs    []*prometheus.Desc
}

// NewExporter returns an initialized exporter
func NewExporter() *Exporter {
	e := &Exporter{
		descs: make([]*prometheus.Desc, len(metricDefs)),
		dial: func(addr string) (LDAPClient, error) {
			conn, err := ldap.DialURL(addr)
			if err != nil {
				return nil, err
			}
			return &ldapClient{conn: conn}, nil
		},
	}
	for i, m := range metricDefs {
		e.descs[i] = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", m.label),
			m.help, nil, nil,
		)
	}
	return e
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, d := range e.descs {
		ch <- d
	}
}

func (e *Exporter) getLDAPConn() (LDAPClient, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.ldapConn != nil {
		return e.ldapConn, nil
	}

	if e.dial == nil {
		return nil, fmt.Errorf("no dial function configured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), ldapTimeout)
	defer cancel()

	resultCh := make(chan struct {
		c   LDAPClient
		err error
	}, 1)

	go func() {
		c, err := e.dial(fmt.Sprintf("ldap://%s:%d", server, port))
		resultCh <- struct {
			c   LDAPClient
			err error
		}{c: c, err: err}
	}()

	select {
	case result := <-resultCh:
		if result.err != nil {
			log.Printf("LDAP connection failed: %v", result.err)
			return nil, result.err
		}
		e.ldapConn = result.c
		return result.c, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("LDAP connection timeout after %v to %s:%d", ldapTimeout, server, port)
	}
}

func (e *Exporter) closeLDAPConn() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.ldapConn != nil {
		_ = e.ldapConn.Close()
		e.ldapConn = nil
	}
}

// Collect reads stats from LDAP connection object into Prometheus objects
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	conn, err := e.getLDAPConn()
	if err != nil {
		log.Printf("Error getting LDAP connection: %v", err)
		return
	}

	data, err := searchLDAP(conn, ldapTimeout)
	if err != nil {
		log.Printf("Error collecting LDAP stats: %v", err)
		e.closeLDAPConn()
		return
	}

	v := reflect.ValueOf(data)
	for i, m := range metricDefs {
		vt := prometheus.GaugeValue
		if m.kind == counterKind {
			vt = prometheus.CounterValue
		}
		ch <- prometheus.MustNewConstMetric(e.descs[i], vt, v.Field(m.fieldIdx).Float())
	}
}
