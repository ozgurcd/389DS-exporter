// Ozgur Demir <ozgurcd@gmail.com>

package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/ozgurcd/389DS-exporter/obj"
)

// Helper function to parse float with error handling
func parseFloatWithDefault(value, fieldName string) float64 {
	if value == "" {
		return 0
	}

	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Printf("[ERROR] %s: %v", fieldName, err)
		return 0
	}
	return result
}

func getStats(server string, port int, timeout time.Duration) (obj.DSData, error) {
	// Create context with configurable timeout for LDAP operations
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Printf("Connecting to LDAP server %s:%d with %v timeout", server, port, timeout)

	// Use a channel to handle timeout for connection
	type connResult struct {
		conn *ldap.Conn
		err  error
	}

	connCh := make(chan connResult, 1)
	go func() {
		conn, err := ldap.DialURL(fmt.Sprintf("ldap://%s:%d", server, port))
		connCh <- connResult{conn: conn, err: err}
	}()

	var conn *ldap.Conn
	select {
	case result := <-connCh:
		if result.err != nil {
			return obj.DSData{}, fmt.Errorf("failed to connect to LDAP server %s:%d: %w", server, port, result.err)
		}
		conn = result.conn
	case <-ctx.Done():
		return obj.DSData{}, fmt.Errorf("LDAP connection timeout after %v to %s:%d", timeout, server, port)
	}
	defer func() { _ = conn.Close() }()

	searchRequest := ldap.NewSearchRequest(
		"cn=monitor",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(objectclass=*)",
		nil,
		nil,
	)

	// Use SearchWithContext for search timeout (if available) or regular Search with timeout handling
	var sr *ldap.SearchResult
	var err error

	searchCh := make(chan struct {
		sr  *ldap.SearchResult
		err error
	}, 1)

	go func() {
		sr, err := conn.Search(searchRequest)
		searchCh <- struct {
			sr  *ldap.SearchResult
			err error
		}{sr: sr, err: err}
	}()

	select {
	case result := <-searchCh:
		sr = result.sr
		err = result.err
	case <-ctx.Done():
		return obj.DSData{}, fmt.Errorf("LDAP search timeout after %v", timeout)
	}

	if err != nil {
		return obj.DSData{}, fmt.Errorf("LDAP search failed: %w", err)
	}

	return parseMonitorAttrs(sr.Entries), nil
}

func parseMonitorAttrs(entries []*ldap.Entry) obj.DSData {
	attrs := make(map[string]string, 33)
	for _, entry := range entries {
		for _, attr := range entry.Attributes {
			if len(attr.Values) == 0 {
				continue
			}
			attrs[attr.Name] = attr.Values[0]
			log.Printf("Found %s: %s", attr.Name, attr.Values[0])
		}
	}

	return obj.DSData{
		Threads:                    parseFloatWithDefault(attrs["threads"], "threads"),
		Readwaiters:                parseFloatWithDefault(attrs["readwaiters"], "readwaiters"),
		Opsinitiated:               parseFloatWithDefault(attrs["opsinitiated"], "opsinitiated"),
		Opscompleted:               parseFloatWithDefault(attrs["opscompleted"], "opscompleted"),
		Dtablesize:                 parseFloatWithDefault(attrs["dtablesize"], "dtablesize"),
		Anonymousbinds:             parseFloatWithDefault(attrs["anonymousbinds"], "anonymousbinds"),
		Unauthbinds:                parseFloatWithDefault(attrs["unauthbinds"], "unauthbinds"),
		Simpleauthbinds:            parseFloatWithDefault(attrs["simpleauthbinds"], "simpleauthbinds"),
		Strongauthbinds:            parseFloatWithDefault(attrs["strongauthbinds"], "strongauthbinds"),
		Bindsecurityerrors:         parseFloatWithDefault(attrs["bindsecurityerrors"], "bindsecurityerrors"),
		Inops:                      parseFloatWithDefault(attrs["inops"], "inops"),
		Readops:                    parseFloatWithDefault(attrs["readops"], "readops"),
		Compareops:                 parseFloatWithDefault(attrs["compareops"], "compareops"),
		Addentryops:                parseFloatWithDefault(attrs["addentryops"], "addentryops"),
		Removeentryops:             parseFloatWithDefault(attrs["removeentryops"], "removeentryops"),
		Modifyentryops:             parseFloatWithDefault(attrs["modifyentryops"], "modifyentryops"),
		Modifyrdnops:               parseFloatWithDefault(attrs["modifyrdnops"], "modifyrdnops"),
		Searchops:                  parseFloatWithDefault(attrs["searchops"], "searchops"),
		Onelevelsearchops:          parseFloatWithDefault(attrs["onelevelsearchops"], "onelevelsearchops"),
		Wholesubtreesearchops:      parseFloatWithDefault(attrs["wholesubtreesearchops"], "wholesubtreesearchops"),
		Referrals:                  parseFloatWithDefault(attrs["referrals"], "referrals"),
		Securityerrors:             parseFloatWithDefault(attrs["securityerrors"], "securityerrors"),
		Errors:                     parseFloatWithDefault(attrs["errors"], "errors"),
		Connections:                parseFloatWithDefault(attrs["connections"], "connections"),
		Connectionseq:              parseFloatWithDefault(attrs["connectionseq"], "connectionseq"),
		Connectionsinmaxthreads:    parseFloatWithDefault(attrs["connectionsinmaxthreads"], "connectionsinmaxthreads"),
		Connectionsmaxthreadscount: parseFloatWithDefault(attrs["connectionsmaxthreadscount"], "connectionsmaxthreadscount"),
		Bytesrecv:                  parseFloatWithDefault(attrs["bytesrecv"], "bytesrecv"),
		Bytessent:                  parseFloatWithDefault(attrs["bytessent"], "bytessent"),
		Entriesreturned:            parseFloatWithDefault(attrs["entriesreturned"], "entriesreturned"),
		Referralsreturned:          parseFloatWithDefault(attrs["referralsreturned"], "referralsreturned"),
		Cacheentries:               parseFloatWithDefault(attrs["cacheentries"], "cacheentries"),
		Cachehits:                  parseFloatWithDefault(attrs["cachehits"], "cachehits"),
	}
}
