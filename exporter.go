package main

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

// Exporter stores metrics from 389DS
type Exporter struct {
	Threads                    *prometheus.Desc
	Readwaiters                *prometheus.Desc
	Opsinitiated               *prometheus.Desc
	Opscompleted               *prometheus.Desc
	Dtablesize                 *prometheus.Desc
	Anonymousbinds             *prometheus.Desc
	Unauthbinds                *prometheus.Desc
	Simpleauthbinds            *prometheus.Desc
	Strongauthbinds            *prometheus.Desc
	Bindsecurityerrors         *prometheus.Desc
	Inops                      *prometheus.Desc
	Readops                    *prometheus.Desc
	Compareops                 *prometheus.Desc
	Addentryops                *prometheus.Desc
	Removeentryops             *prometheus.Desc
	Modifyentryops             *prometheus.Desc
	Modifyrdnops               *prometheus.Desc
	Searchops                  *prometheus.Desc
	Onelevelsearchops          *prometheus.Desc
	Wholesubtreesearchops      *prometheus.Desc
	Referrals                  *prometheus.Desc
	Securityerrors             *prometheus.Desc
	Errors                     *prometheus.Desc
	Connections                *prometheus.Desc
	Connectionseq              *prometheus.Desc
	Connectionsinmaxthreads    *prometheus.Desc
	Connectionsmaxthreadscount *prometheus.Desc
	Bytesrecv                  *prometheus.Desc
	Bytessent                  *prometheus.Desc
	Entriesreturned            *prometheus.Desc
	Referralsreturned          *prometheus.Desc
	Cacheentries               *prometheus.Desc
	Cachehits                  *prometheus.Desc
}

// NewExporter returns an initialized exporter
func NewExporter() *Exporter {
	return &Exporter{

		Threads: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "threads"),
			"Number of Threads max configured",
			nil,
			nil,
		),

		Readwaiters: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "readwaiters"),
			"Current number of threads waiting to read data from a client",
			nil,
			nil,
		),

		Opsinitiated: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "opsinitiated"),
			"Current number of operations the server has initiated since it started",
			nil,
			nil,
		),

		Opscompleted: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "opscompleted"),
			"Current number of operations the server has completed since it started",
			nil,
			nil,
		),

		Dtablesize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "dtablesize"),
			"The number of file descriptors available to the directory. Essentially, this value shows how many additional concurrent connections can be serviced by the directory",
			nil,
			nil,
		),

		Anonymousbinds: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "anonymousbinds"),
			"Number of Anonymous Binds",
			nil,
			nil,
		),

		Unauthbinds: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "unauthbinds"),
			"Number of Unauth Binds",
			nil,
			nil,
		),

		Simpleauthbinds: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "simpleauthbinds"),
			"Number of Simple Auth Binds",
			nil,
			nil,
		),

		Strongauthbinds: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "strongauthbinds"),
			"Number of Strong Auth Binds",
			nil,
			nil,
		),

		Bindsecurityerrors: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "bindsecurityerrors"),
			"Number of Bind Security Errors",
			nil,
			nil,
		),

		Inops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "inops"),
			"Number of All Requests",
			nil,
			nil,
		),

		Readops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "readops"),
			"Number of Read Operations",
			nil,
			nil,
		),

		Compareops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "compareops"),
			"Number of Compare Operations",
			nil,
			nil,
		),

		Addentryops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "addentryops"),
			"Number of Add Entry Operations",
			nil,
			nil,
		),

		Removeentryops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "removeentryops"),
			"Number of Remove Entry Operations",
			nil,
			nil,
		),

		Modifyentryops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "modifyentryops"),
			"Number of Modify Entry Operations",
			nil,
			nil,
		),

		Modifyrdnops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "modifyrdnops"),
			"Number of Modify RDN Operations",
			nil,
			nil,
		),

		Searchops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "searchops"),
			"Number of LDAP Search Requests",
			nil,
			nil,
		),

		Onelevelsearchops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "onelevelsearchops"),
			"Number of one-level Search Requests",
			nil,
			nil,
		),

		Wholesubtreesearchops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "wholesubtreesearchops"),
			"Number of subtree-level Search Requests",
			nil,
			nil,
		),

		Referrals: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "referrals"),
			"Number of LDAP referrals",
			nil,
			nil,
		),

		Securityerrors: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "securityerrors"),
			"Number of Security Errors",
			nil,
			nil,
		),

		Errors: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "errors"),
			"Number of Errors",
			nil,
			nil,
		),

		Connections: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "connections"),
			"Number of Connections in Open State at the sampling time",
			nil,
			nil,
		),

		Connectionseq: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "connectionseq"),
			"Total Number of Connections opened",
			nil,
			nil,
		),

		Connectionsinmaxthreads: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "connectionsinmaxthreads"),
			"Number of connections that are currently in a max thread state",
			nil,
			nil,
		),

		Connectionsmaxthreadscount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "connectionsmaxthreadscount"),
			"Number of connectionsmaxthreadscount",
			nil,
			nil,
		),

		Bytesrecv: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "bytesrecv"),
			"Total number of bytes received",
			nil,
			nil,
		),

		Bytessent: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "bytessent"),
			"Total number of bytes sent",
			nil,
			nil,
		),

		Entriesreturned: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "entriesreturned"),
			"Number of Entries Returned",
			nil,
			nil,
		),

		Referralsreturned: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "referralsreturned"),
			"Number of Referrals Returned",
			nil,
			nil,
		),

		Cacheentries: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "cacheentries"),
			"Number of Cache Entries",
			nil,
			nil,
		),

		Cachehits: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "cachehits"),
			"Number of Cache Hits",
			nil,
			nil,
		),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.Threads
	ch <- e.Readwaiters
	ch <- e.Opsinitiated
	ch <- e.Opscompleted
	ch <- e.Dtablesize
	ch <- e.Anonymousbinds
	ch <- e.Unauthbinds
	ch <- e.Simpleauthbinds
	ch <- e.Strongauthbinds
	ch <- e.Bindsecurityerrors
	ch <- e.Inops
	ch <- e.Readops
	ch <- e.Compareops
	ch <- e.Addentryops
	ch <- e.Removeentryops
	ch <- e.Modifyentryops
	ch <- e.Modifyrdnops
	ch <- e.Searchops
	ch <- e.Onelevelsearchops
	ch <- e.Wholesubtreesearchops
	ch <- e.Referrals
	ch <- e.Securityerrors
	ch <- e.Errors
	ch <- e.Connections
	ch <- e.Connectionseq
	ch <- e.Connectionsinmaxthreads
	ch <- e.Connectionsmaxthreadscount
	ch <- e.Bytesrecv
	ch <- e.Bytessent
	ch <- e.Entriesreturned
	ch <- e.Referralsreturned
	ch <- e.Cacheentries
	ch <- e.Cachehits
}

// Collect reads stats from LDAP connection object into Prometheus objects
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	data, err := getStats(server, port, ldapTimeout)
	if err != nil {
		log.Printf("Error collecting LDAP stats: %v", err)
		// Return early - no metrics will be exported if LDAP is unavailable
		return
	}

	// Current state metrics - use GaugeValue
	ch <- prometheus.MustNewConstMetric(e.Threads, prometheus.GaugeValue, data.Threads)
	ch <- prometheus.MustNewConstMetric(e.Readwaiters, prometheus.GaugeValue, data.Readwaiters)
	ch <- prometheus.MustNewConstMetric(e.Dtablesize, prometheus.GaugeValue, data.Dtablesize)
	ch <- prometheus.MustNewConstMetric(e.Connections, prometheus.GaugeValue, data.Connections)
	ch <- prometheus.MustNewConstMetric(e.Connectionsinmaxthreads, prometheus.GaugeValue, data.Connectionsinmaxthreads)
	ch <- prometheus.MustNewConstMetric(e.Connectionsmaxthreadscount, prometheus.GaugeValue, data.Connectionsmaxthreadscount)
	ch <- prometheus.MustNewConstMetric(e.Cacheentries, prometheus.GaugeValue, data.Cacheentries)

	// Cumulative counters - use CounterValue (only ever increase)
	ch <- prometheus.MustNewConstMetric(e.Opsinitiated, prometheus.CounterValue, data.Opsinitiated)
	ch <- prometheus.MustNewConstMetric(e.Opscompleted, prometheus.CounterValue, data.Opscompleted)
	ch <- prometheus.MustNewConstMetric(e.Anonymousbinds, prometheus.CounterValue, data.Anonymousbinds)
	ch <- prometheus.MustNewConstMetric(e.Unauthbinds, prometheus.CounterValue, data.Unauthbinds)
	ch <- prometheus.MustNewConstMetric(e.Simpleauthbinds, prometheus.CounterValue, data.Simpleauthbinds)
	ch <- prometheus.MustNewConstMetric(e.Strongauthbinds, prometheus.CounterValue, data.Strongauthbinds)
	ch <- prometheus.MustNewConstMetric(e.Bindsecurityerrors, prometheus.CounterValue, data.Bindsecurityerrors)
	ch <- prometheus.MustNewConstMetric(e.Inops, prometheus.CounterValue, data.Inops)
	ch <- prometheus.MustNewConstMetric(e.Readops, prometheus.CounterValue, data.Readops)
	ch <- prometheus.MustNewConstMetric(e.Compareops, prometheus.CounterValue, data.Compareops)
	ch <- prometheus.MustNewConstMetric(e.Addentryops, prometheus.CounterValue, data.Addentryops)
	ch <- prometheus.MustNewConstMetric(e.Removeentryops, prometheus.CounterValue, data.Removeentryops)
	ch <- prometheus.MustNewConstMetric(e.Modifyentryops, prometheus.CounterValue, data.Modifyentryops)
	ch <- prometheus.MustNewConstMetric(e.Modifyrdnops, prometheus.CounterValue, data.Modifyrdnops)
	ch <- prometheus.MustNewConstMetric(e.Searchops, prometheus.CounterValue, data.Searchops)
	ch <- prometheus.MustNewConstMetric(e.Onelevelsearchops, prometheus.CounterValue, data.Onelevelsearchops)
	ch <- prometheus.MustNewConstMetric(e.Wholesubtreesearchops, prometheus.CounterValue, data.Wholesubtreesearchops)
	ch <- prometheus.MustNewConstMetric(e.Referrals, prometheus.CounterValue, data.Referrals)
	ch <- prometheus.MustNewConstMetric(e.Securityerrors, prometheus.CounterValue, data.Securityerrors)
	ch <- prometheus.MustNewConstMetric(e.Errors, prometheus.CounterValue, data.Errors)
	ch <- prometheus.MustNewConstMetric(e.Connectionseq, prometheus.CounterValue, data.Connectionseq)
	ch <- prometheus.MustNewConstMetric(e.Bytesrecv, prometheus.CounterValue, data.Bytesrecv)
	ch <- prometheus.MustNewConstMetric(e.Bytessent, prometheus.CounterValue, data.Bytessent)
	ch <- prometheus.MustNewConstMetric(e.Entriesreturned, prometheus.CounterValue, data.Entriesreturned)
	ch <- prometheus.MustNewConstMetric(e.Referralsreturned, prometheus.CounterValue, data.Referralsreturned)
	ch <- prometheus.MustNewConstMetric(e.Cachehits, prometheus.CounterValue, data.Cachehits)
}
