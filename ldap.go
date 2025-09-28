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
	defer conn.Close()

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

	var (
		threads                    string
		readwaiters                string
		opsinitiated               string
		opscompleted               string
		dtablesize                 string
		anonymousbinds             string
		unauthbinds                string
		simpleauthbinds            string
		strongauthbinds            string
		bindsecurityerrors         string
		inops                      string
		readops                    string
		compareops                 string
		addentryops                string
		removeentryops             string
		modifyentryops             string
		modifyrdnops               string
		searchops                  string
		onelevelsearchops          string
		wholesubtreesearchops      string
		referrals                  string
		securityerrors             string
		errors                     string
		connections                string
		connectionseq              string
		connectionsinmaxthreads    string
		connectionsmaxthreadscount string
		bytesrecv                  string
		bytessent                  string
		entriesreturned            string
		referralsreturned          string
		cacheentries               string
		cachehits                  string
	)

	for n := 0; n < len(sr.Entries); n++ {
		entry := sr.Entries[n]
		attributes := entry.Attributes

		for _, attr := range attributes {
			name := string(attr.Name)
			value := attr.Values

			// Only log if we're interested in this attribute and it has values
			if len(value) > 0 {
				switch name {
				case "threads":
					threads = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "readwaiters":
					readwaiters = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "opsinitiated":
					opsinitiated = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "opscompleted":
					opscompleted = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "dtablesize":
					dtablesize = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "anonymousbinds":
					anonymousbinds = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "unauthbinds":
					unauthbinds = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "simpleauthbinds":
					simpleauthbinds = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "strongauthbinds":
					strongauthbinds = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "bindsecurityerrors":
					bindsecurityerrors = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "inops":
					inops = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "readops":
					readops = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "compareops":
					compareops = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "addentryops":
					addentryops = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "removeentryops":
					removeentryops = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "modifyentryops":
					modifyentryops = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "modifyrdnops":
					modifyrdnops = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "searchops":
					searchops = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "onelevelsearchops":
					onelevelsearchops = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "wholesubtreesearchops":
					wholesubtreesearchops = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "referrals":
					referrals = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "securityerrors":
					securityerrors = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "errors":
					errors = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "connections":
					connections = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "connectionseq":
					connectionseq = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "connectionsinmaxthreads":
					connectionsinmaxthreads = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "connectionsmaxthreadscount":
					connectionsmaxthreadscount = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "bytesrecv":
					bytesrecv = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "bytessent":
					bytessent = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "entriesreturned":
					entriesreturned = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "referralsreturned":
					referralsreturned = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "cacheentries":
					cacheentries = value[0]
					log.Printf("Found %s: %s", name, value[0])
				case "cachehits":
					cachehits = value[0]
					log.Printf("Found %s: %s", name, value[0])
				}
			}
		}
	}

	return obj.DSData{
		Threads:                    parseFloatWithDefault(threads, "threads"),
		Readwaiters:                parseFloatWithDefault(readwaiters, "readwaiters"),
		Opsinitiated:               parseFloatWithDefault(opsinitiated, "opsinitiated"),
		Opscompleted:               parseFloatWithDefault(opscompleted, "opscompleted"),
		Dtablesize:                 parseFloatWithDefault(dtablesize, "dtablesize"),
		Anonymousbinds:             parseFloatWithDefault(anonymousbinds, "anonymousbinds"),
		Unauthbinds:                parseFloatWithDefault(unauthbinds, "unauthbinds"),
		Simpleauthbinds:            parseFloatWithDefault(simpleauthbinds, "simpleauthbinds"),
		Strongauthbinds:            parseFloatWithDefault(strongauthbinds, "strongauthbinds"),
		Bindsecurityerrors:         parseFloatWithDefault(bindsecurityerrors, "bindsecurityerrors"),
		Inops:                      parseFloatWithDefault(inops, "inops"),
		Readops:                    parseFloatWithDefault(readops, "readops"),
		Compareops:                 parseFloatWithDefault(compareops, "compareops"),
		Addentryops:                parseFloatWithDefault(addentryops, "addentryops"),
		Removeentryops:             parseFloatWithDefault(removeentryops, "removeentryops"),
		Modifyentryops:             parseFloatWithDefault(modifyentryops, "modifyentryops"),
		Modifyrdnops:               parseFloatWithDefault(modifyrdnops, "modifyrdnops"),
		Searchops:                  parseFloatWithDefault(searchops, "searchops"),
		Onelevelsearchops:          parseFloatWithDefault(onelevelsearchops, "onelevelsearchops"),
		Wholesubtreesearchops:      parseFloatWithDefault(wholesubtreesearchops, "wholesubtreesearchops"),
		Referrals:                  parseFloatWithDefault(referrals, "referrals"),
		Securityerrors:             parseFloatWithDefault(securityerrors, "securityerrors"),
		Errors:                     parseFloatWithDefault(errors, "errors"),
		Connections:                parseFloatWithDefault(connections, "connections"),
		Connectionseq:              parseFloatWithDefault(connectionseq, "connectionseq"),
		Connectionsinmaxthreads:    parseFloatWithDefault(connectionsinmaxthreads, "connectionsinmaxthreads"),
		Connectionsmaxthreadscount: parseFloatWithDefault(connectionsmaxthreadscount, "connectionsmaxthreadscount"),
		Bytesrecv:                  parseFloatWithDefault(bytesrecv, "bytesrecv"),
		Bytessent:                  parseFloatWithDefault(bytessent, "bytessent"),
		Entriesreturned:            parseFloatWithDefault(entriesreturned, "entriesreturned"),
		Referralsreturned:          parseFloatWithDefault(referralsreturned, "referralsreturned"),
		Cacheentries:               parseFloatWithDefault(cacheentries, "cacheentries"),
		Cachehits:                  parseFloatWithDefault(cachehits, "cachehits"),
	}, nil
}
