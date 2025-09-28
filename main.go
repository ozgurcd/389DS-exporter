// Ozgur Demir <ozgurcd@gmail.com>

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/spf13/pflag"
)

const (
	namespace = "ds_exporter"
)

var (
	port        int
	server      string
	_version    = "1.6"
	ldapTimeout time.Duration
)

func main() {
	var (
		listenAddress  = pflag.String("web.listen-address", ":9313", "Address to listen on for web interface and telemetry.")
		metricsPath    = pflag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
		ldapServer     = pflag.String("ldap.ServerFQDN", "localhost", "FQDN of the target LDAP server")
		ldapServerPort = pflag.Int("ldap.ServerPort", 389, "Port to connect on LDAP server")
		timeout        = pflag.Duration("ldap.timeout", 10*time.Second, "LDAP connection timeout")
		showVersion    = pflag.BoolP("version", "v", false, "Show version information")
		showHelp       = pflag.BoolP("help", "h", false, "Show help")
	)

	pflag.Parse()

	if *showHelp {
		pflag.Usage()
		return
	}

	if *showVersion {
		println(version.Print("ds_exporter"))
		return
	}

	// Validate configuration
	if *ldapServerPort < 1 || *ldapServerPort > 65535 {
		log.Fatal("Invalid LDAP port number: must be between 1 and 65535")
	}

	if *ldapServer == "" {
		log.Fatal("LDAP server cannot be empty")
	}

	port = *ldapServerPort
	server = *ldapServer
	ldapTimeout = *timeout
	version.Version = _version

	log.Println("Starting ds_exporter", version.Info())
	log.Println("Build context", version.BuildContext())
	log.Printf("Target LDAP Server: %s:%d (timeout: %v)", *ldapServer, port, *timeout)

	prometheus.MustRegister(NewExporter())

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>389-DS Exporter</title></head>
             <body>
             <h1>389-DS Exporter</h1>
             <p>For the metrics: Click <a href='` + *metricsPath + `'>here</a></p>
             <p>Health check: <a href='/health'>here</a></p>
             </body>
             </html>`))
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, err := getStats(server, port, ldapTimeout)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("LDAP connection failed: " + err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Create HTTP server with timeouts
	srv := &http.Server{
		Addr:         *listenAddress,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in background
	go func() {
		log.Println("Starting HTTP server on", *listenAddress)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal("HTTP server failed:", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped")
}
