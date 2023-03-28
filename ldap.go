// Ozgur Demir <ozgurcd@gmail.com>

package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-ldap/ldap"
)

func getStats(server string, port int) DSData {
	log.Println("ldap server:", server)
	log.Println("ldap port:", port)
	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", server, port))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	searchRequest := ldap.NewSearchRequest(
		"cn=monitor",
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
			for i := 0; i < len(value); i++ {
				log.Printf("%v",
					fmt.Sprintf("Name: %s, Value: %s[%d]\n",
						name,
						value[i], i))
			}

			// ignoring unused attributes.
			switch name {
			case "threads":
				threads = value[0]
			case "readwaiters":
				readwaiters = value[0]
			case "opsinitiated":
				opsinitiated = value[0]
			case "opscompleted":
				opscompleted = value[0]
			case "dtablesize":
				dtablesize = value[0]
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
			}
		}
	}

	threads64, err := strconv.ParseFloat(threads, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]threads %v\n", err)
		}
		threads64 = 0
	}

	readwaiters64, err := strconv.ParseFloat(readwaiters, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]readwaiters %v\n", err)
		}
		readwaiters64 = 0
	}

	opsinitiated64, err := strconv.ParseFloat(opsinitiated, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]opsinitiated %v\n", err)
		}
		opsinitiated64 = 0
	}

	opscompleted64, err := strconv.ParseFloat(opscompleted, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]opscompleted %v\n", err)
		}
		opscompleted64 = 0
	}

	dtablesize64, err := strconv.ParseFloat(dtablesize, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]dtablesize %v\n", err)
		}
		dtablesize64 = 0
	}

	anonymousbinds64, err := strconv.ParseFloat(anonymousbinds, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]anonymousbinds %v\n", err)
		}
		anonymousbinds64 = 0
	}

	unauthbinds64, err := strconv.ParseFloat(unauthbinds, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]unauthbinds %v\n", err)
		}
		unauthbinds64 = 0
	}

	simpleauthbinds64, err := strconv.ParseFloat(simpleauthbinds, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]simpleauthbinds %v\n", err)
		}
		simpleauthbinds64 = 0
	}

	strongauthbinds64, err := strconv.ParseFloat(strongauthbinds, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]strongauthbinds %v\n", err)
		}
		strongauthbinds64 = 0
	}

	bindsecurityerrors64, err := strconv.ParseFloat(bindsecurityerrors, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]bindsecurityerrors %v\n", err)
		}
		bindsecurityerrors64 = 0
	}

	inops64, err := strconv.ParseFloat(inops, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]inops %v\n", err)
		}
		inops64 = 0
	}

	readops64, err := strconv.ParseFloat(readops, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]readops %v\n", err)
		}
		readops64 = 0
	}

	compareops64, err := strconv.ParseFloat(compareops, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]compareops %v\n", err)
		}
		compareops64 = 0
	}

	addentryops64, err := strconv.ParseFloat(addentryops, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]addentryops %v\n", err)
		}
		addentryops64 = 0
	}

	removeentryops64, err := strconv.ParseFloat(removeentryops, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]removeentryops %v\n", err)
		}
		removeentryops64 = 0
	}

	modifyentryops64, err := strconv.ParseFloat(modifyentryops, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]modifyentryops %v\n", err)
		}
		modifyentryops64 = 0
	}

	modifyrdnops64, err := strconv.ParseFloat(modifyrdnops, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]modifyrdnops %v\n", err)
		}
		modifyrdnops64 = 0
	}

	searchops64, err := strconv.ParseFloat(searchops, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]searchops %v\n", err)
		}
		searchops64 = 0
	}

	onelevelsearchops64, err := strconv.ParseFloat(onelevelsearchops, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]onelevelsearchops %v\n", err)
		}
		onelevelsearchops64 = 0
	}

	wholesubtreesearchops64, err := strconv.ParseFloat(wholesubtreesearchops, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]wholesubtreesearchops %v\n", err)
		}
		wholesubtreesearchops64 = 0
	}

	referrals64, err := strconv.ParseFloat(referrals, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]referrals %v\n", err)
		}
		referrals64 = 0
	}

	securityerrors64, err := strconv.ParseFloat(securityerrors, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]securityerrors %v\n", err)
		}
		securityerrors64 = 0
	}

	errors64, err := strconv.ParseFloat(errors, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]errors %v\n", err)
		}
		errors64 = 0
	}

	connections64, err := strconv.ParseFloat(connections, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]connections %v\n", err)
		}
		connections64 = 0
	}

	connectionseq64, err := strconv.ParseFloat(connectionseq, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]connectionseq %v\n", err)
		}
		connectionseq64 = 0
	}

	connectionsinmaxthreads64, err := strconv.ParseFloat(connectionsinmaxthreads, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]connectionsinmaxthreads %v\n", err)
		}
		connectionsinmaxthreads64 = 0
	}

	connectionsmaxthreadscount64, err := strconv.ParseFloat(connectionsmaxthreadscount, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]connectionsmaxthreadscount %v\n", err)
		}
		connectionsmaxthreadscount64 = 0
	}

	bytesrecv64, err := strconv.ParseFloat(bytesrecv, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]bytesrecv %v\n", err)
		}
		bytesrecv64 = 0
	}

	bytessent64, err := strconv.ParseFloat(bytessent, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]bytessent %v\n", err)
		}
		bytessent64 = 0
	}

	entriesreturned64, err := strconv.ParseFloat(entriesreturned, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]entriesreturned %v\n", err)
		}
		entriesreturned64 = 0
	}

	referralsreturned64, err := strconv.ParseFloat(referralsreturned, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]referralsreturned %v\n", err)
		}
		referralsreturned64 = 0
	}

	cacheentries64, err := strconv.ParseFloat(cacheentries, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]cacheentries %v\n", err)
		}
		cacheentries64 = 0
	}

	cachehits64, err := strconv.ParseFloat(cachehits, 64)
	if err != nil {
		if errors != "" {
			log.Printf("[ERROR]cachehits %v\n", err)
		}
		cachehits64 = 0
	}

	return DSData{
		threads64,
		readwaiters64,
		opsinitiated64,
		opscompleted64,
		dtablesize64,
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
