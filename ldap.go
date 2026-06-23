// Ozgur Demir <ozgurcd@gmail.com>

package main

import (
	"fmt"
	"log"
	"reflect"
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

func searchLDAP(conn LDAPClient, timeout time.Duration) (obj.DSData, error) {
	searchRequest := ldap.NewSearchRequest(
		"cn=monitor",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(objectclass=*)",
		nil,
		nil,
	)

	return runWithTimeout(timeout, func() (obj.DSData, error) {
		sr, err := conn.Search(searchRequest)
		if err != nil {
			return obj.DSData{}, fmt.Errorf("LDAP search failed: %w", err)
		}
		if sr == nil {
			return obj.DSData{}, fmt.Errorf("LDAP search returned nil result")
		}
		return parseMonitorAttrs(sr.Entries), nil
	})
}

func parseMonitorAttrs(entries []*ldap.Entry) obj.DSData {
	var d obj.DSData
	v := reflect.ValueOf(&d).Elem()
	for _, entry := range entries {
		for _, attr := range entry.Attributes {
			if len(attr.Values) == 0 {
				continue
			}
			if idx, ok := ldapFieldMap[attr.Name]; ok {
				v.Field(idx).SetFloat(parseFloatWithDefault(attr.Values[0], attr.Name))
			}
		}
	}
	return d
}
