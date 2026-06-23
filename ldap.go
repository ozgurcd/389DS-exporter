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

func searchLDAP(conn *ldap.Conn, timeout time.Duration) (obj.DSData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	searchRequest := ldap.NewSearchRequest(
		"cn=monitor",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(objectclass=*)",
		nil,
		nil,
	)

	type searchResult struct {
		sr  *ldap.SearchResult
		err error
	}

	resultCh := make(chan searchResult, 1)
	go func() {
		sr, err := conn.Search(searchRequest)
		resultCh <- searchResult{sr: sr, err: err}
	}()

	select {
	case result := <-resultCh:
		if result.err != nil {
			return obj.DSData{}, fmt.Errorf("LDAP search failed: %w", result.err)
		}
		return parseMonitorAttrs(result.sr.Entries), nil
	case <-ctx.Done():
		return obj.DSData{}, fmt.Errorf("LDAP search timeout after %v", timeout)
	}
}

func parseMonitorAttrs(entries []*ldap.Entry) obj.DSData {
	var d obj.DSData
	for _, entry := range entries {
		for _, attr := range entry.Attributes {
			if len(attr.Values) == 0 {
				continue
			}
			val := attr.Values[0]
			switch attr.Name {
			case "threads":
				d.Threads = parseFloatWithDefault(val, "threads")
			case "readwaiters":
				d.Readwaiters = parseFloatWithDefault(val, "readwaiters")
			case "opsinitiated":
				d.Opsinitiated = parseFloatWithDefault(val, "opsinitiated")
			case "opscompleted":
				d.Opscompleted = parseFloatWithDefault(val, "opscompleted")
			case "dtablesize":
				d.Dtablesize = parseFloatWithDefault(val, "dtablesize")
			case "anonymousbinds":
				d.Anonymousbinds = parseFloatWithDefault(val, "anonymousbinds")
			case "unauthbinds":
				d.Unauthbinds = parseFloatWithDefault(val, "unauthbinds")
			case "simpleauthbinds":
				d.Simpleauthbinds = parseFloatWithDefault(val, "simpleauthbinds")
			case "strongauthbinds":
				d.Strongauthbinds = parseFloatWithDefault(val, "strongauthbinds")
			case "bindsecurityerrors":
				d.Bindsecurityerrors = parseFloatWithDefault(val, "bindsecurityerrors")
			case "inops":
				d.Inops = parseFloatWithDefault(val, "inops")
			case "readops":
				d.Readops = parseFloatWithDefault(val, "readops")
			case "compareops":
				d.Compareops = parseFloatWithDefault(val, "compareops")
			case "addentryops":
				d.Addentryops = parseFloatWithDefault(val, "addentryops")
			case "removeentryops":
				d.Removeentryops = parseFloatWithDefault(val, "removeentryops")
			case "modifyentryops":
				d.Modifyentryops = parseFloatWithDefault(val, "modifyentryops")
			case "modifyrdnops":
				d.Modifyrdnops = parseFloatWithDefault(val, "modifyrdnops")
			case "searchops":
				d.Searchops = parseFloatWithDefault(val, "searchops")
			case "onelevelsearchops":
				d.Onelevelsearchops = parseFloatWithDefault(val, "onelevelsearchops")
			case "wholesubtreesearchops":
				d.Wholesubtreesearchops = parseFloatWithDefault(val, "wholesubtreesearchops")
			case "referrals":
				d.Referrals = parseFloatWithDefault(val, "referrals")
			case "securityerrors":
				d.Securityerrors = parseFloatWithDefault(val, "securityerrors")
			case "errors":
				d.Errors = parseFloatWithDefault(val, "errors")
			case "connections":
				d.Connections = parseFloatWithDefault(val, "connections")
			case "connectionseq":
				d.Connectionseq = parseFloatWithDefault(val, "connectionseq")
			case "connectionsinmaxthreads":
				d.Connectionsinmaxthreads = parseFloatWithDefault(val, "connectionsinmaxthreads")
			case "connectionsmaxthreadscount":
				d.Connectionsmaxthreadscount = parseFloatWithDefault(val, "connectionsmaxthreadscount")
			case "bytesrecv":
				d.Bytesrecv = parseFloatWithDefault(val, "bytesrecv")
			case "bytessent":
				d.Bytessent = parseFloatWithDefault(val, "bytessent")
			case "entriesreturned":
				d.Entriesreturned = parseFloatWithDefault(val, "entriesreturned")
			case "referralsreturned":
				d.Referralsreturned = parseFloatWithDefault(val, "referralsreturned")
			case "cacheentries":
				d.Cacheentries = parseFloatWithDefault(val, "cacheentries")
			case "cachehits":
				d.Cachehits = parseFloatWithDefault(val, "cachehits")
			}
		}
	}
	return d
}
