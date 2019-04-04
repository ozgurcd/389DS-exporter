// Ozgur Demir <ozgurcd@gmail.com>

package main

import (
	"fmt"
	"strconv"

	"github.com/prometheus/common/log"
	ldap "gopkg.in/ldap.v3"
)

func getStats(ldapServer *string, ldapServerPort *string) DSData {
	port, err := strconv.Atoi(*ldapServerPort)
	if nil != err {
		log.Fatal(err)
	}

	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", *ldapServer, port))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	searchRequest := ldap.NewSearchRequest(
		"cn=snmp, cn=monitor",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(objectclass=*)",
		nil,
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	var (
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

	entry := sr.Entries[0]
	attributes := entry.Attributes

	for _, attr := range attributes {
		name := string(attr.Name)
		value := attr.Values

		// ignoring unused attributes.
		switch name {
		case "anonymousbinds":
			anonymousbinds = value[0]
		case "unauthbinds":
			unauthbinds = value[0]
		case "simpleauthbinds":
			simpleauthbinds = value[0]
		case "strongauthbinds":
			strongauthbinds = value[0]
		case "bindsecurityerrors":
			bindsecurityerrors = value[0]
		case "inops":
			inops = value[0]
		case "readops":
			readops = value[0]
		case "compareops":
			compareops = value[0]
		case "addentryops":
			addentryops = value[0]
		case "removeentryops":
			removeentryops = value[0]
		case "modifyentryops":
			modifyentryops = value[0]
		case "modifyrdnops":
			modifyrdnops = value[0]
		case "searchops":
			searchops = value[0]
		case "onelevelsearchops":
			onelevelsearchops = value[0]
		case "wholesubtreesearchops":
			wholesubtreesearchops = value[0]
		case "referrals":
			referrals = value[0]
		case "securityerrors":
			securityerrors = value[0]
		case "errors":
			errors = value[0]
		case "connections":
			connections = value[0]
		case "connectionseq":
			connectionseq = value[0]
		case "connectionsinmaxthreads":
			connectionsinmaxthreads = value[0]
		case "connectionsmaxthreadscount":
			connectionsmaxthreadscount = value[0]
		case "bytesrecv":
			bytesrecv = value[0]
		case "bytessent":
			bytessent = value[0]
		case "entriesreturned":
			entriesreturned = value[0]
		case "referralsreturned":
			referralsreturned = value[0]
		case "cacheentries":
			cacheentries = value[0]
		case "cachehits":
			cachehits = value[0]
		default:
			//fmt.Printf("Name: %s, Value: %s\n", name, value)
		}
	}

	anonymousbinds64, err := strconv.ParseFloat(anonymousbinds, 64)
	if err != nil {
		log.Error(err)
		anonymousbinds64 = 0
	}

	unauthbinds64, err := strconv.ParseFloat(unauthbinds, 64)
	if err != nil {
		log.Error(err)
		unauthbinds64 = 0
	}

	simpleauthbinds64, err := strconv.ParseFloat(simpleauthbinds, 64)
	if err != nil {
		log.Error(err)
		simpleauthbinds64 = 0
	}

	strongauthbinds64, err := strconv.ParseFloat(strongauthbinds, 64)
	if err != nil {
		log.Error(err)
		strongauthbinds64 = 0
	}

	bindsecurityerrors64, err := strconv.ParseFloat(bindsecurityerrors, 64)
	if err != nil {
		log.Error(err)
		bindsecurityerrors64 = 0
	}

	inops64, err := strconv.ParseFloat(inops, 64)
	if err != nil {
		log.Error(err)
		inops64 = 0
	}

	readops64, err := strconv.ParseFloat(readops, 64)
	if err != nil {
		log.Error(err)
		readops64 = 0
	}

	compareops64, err := strconv.ParseFloat(compareops, 64)
	if err != nil {
		log.Error(err)
		compareops64 = 0
	}

	addentryops64, err := strconv.ParseFloat(addentryops, 64)
	if err != nil {
		log.Error(err)
		addentryops64 = 0
	}

	removeentryops64, err := strconv.ParseFloat(removeentryops, 64)
	if err != nil {
		log.Error(err)
		removeentryops64 = 0
	}

	modifyentryops64, err := strconv.ParseFloat(modifyentryops, 64)
	if err != nil {
		log.Error(err)
		modifyentryops64 = 0
	}

	modifyrdnops64, err := strconv.ParseFloat(modifyrdnops, 64)
	if err != nil {
		log.Error(err)
		modifyrdnops64 = 0
	}

	searchops64, err := strconv.ParseFloat(searchops, 64)
	if err != nil {
		log.Error(err)
		searchops64 = 0
	}

	onelevelsearchops64, err := strconv.ParseFloat(onelevelsearchops, 64)
	if err != nil {
		log.Error(err)
		onelevelsearchops64 = 0
	}

	wholesubtreesearchops64, err := strconv.ParseFloat(wholesubtreesearchops, 64)
	if err != nil {
		log.Error(err)
		wholesubtreesearchops64 = 0
	}

	referrals64, err := strconv.ParseFloat(referrals, 64)
	if err != nil {
		log.Error(err)
		referrals64 = 0
	}

	securityerrors64, err := strconv.ParseFloat(securityerrors, 64)
	if err != nil {
		log.Error(err)
		securityerrors64 = 0
	}

	errors64, err := strconv.ParseFloat(errors, 64)
	if err != nil {
		log.Error(err)
		errors64 = 0
	}

	connections64, err := strconv.ParseFloat(connections, 64)
	if err != nil {
		log.Error(err)
		connections64 = 0
	}

	connectionseq64, err := strconv.ParseFloat(connectionseq, 64)
	if err != nil {
		log.Error(err)
		connectionseq64 = 0
	}

	connectionsinmaxthreads64, err := strconv.ParseFloat(connectionsinmaxthreads, 64)
	if err != nil {
		log.Error(err)
		connectionsinmaxthreads64 = 0
	}

	connectionsmaxthreadscount64, err := strconv.ParseFloat(connectionsmaxthreadscount, 64)
	if err != nil {
		log.Error(err)
		connectionsmaxthreadscount64 = 0
	}

	bytesrecv64, err := strconv.ParseFloat(bytesrecv, 64)
	if err != nil {
		log.Error(err)
		bytesrecv64 = 0
	}

	bytessent64, err := strconv.ParseFloat(bytessent, 64)
	if err != nil {
		log.Error(err)
		bytessent64 = 0
	}

	entriesreturned64, err := strconv.ParseFloat(entriesreturned, 64)
	if err != nil {
		log.Error(err)
		entriesreturned64 = 0
	}

	referralsreturned64, err := strconv.ParseFloat(referralsreturned, 64)
	if err != nil {
		log.Error(err)
		referralsreturned64 = 0
	}

	cacheentries64, err := strconv.ParseFloat(cacheentries, 64)
	if err != nil {
		log.Error(err)
		cacheentries64 = 0
	}

	cachehits64, err := strconv.ParseFloat(cachehits, 64)
	if err != nil {
		log.Error(err)
		cachehits64 = 0
	}

	return DSData{
		anonymousbinds64,
		unauthbinds64,
		simpleauthbinds64,
		strongauthbinds64,
		bindsecurityerrors64,
		inops64,
		readops64,
		compareops64,
		addentryops64,
		removeentryops64,
		modifyentryops64,
		modifyrdnops64,
		searchops64,
		onelevelsearchops64,
		wholesubtreesearchops64,
		referrals64,
		securityerrors64,
		errors64,
		connections64,
		connectionseq64,
		connectionsinmaxthreads64,
		connectionsmaxthreadscount64,
		bytesrecv64,
		bytessent64,
		entriesreturned64,
		referralsreturned64,
		cacheentries64,
		cachehits64}
}
