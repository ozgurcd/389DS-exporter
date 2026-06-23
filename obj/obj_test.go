package obj

import (
	"reflect"
	"testing"
)

func TestDSDataZeroValues(t *testing.T) {
	var d DSData
	v := reflect.ValueOf(d)
	for i := range v.NumField() {
		if v.Field(i).Float() != 0 {
			t.Errorf("%s should default to 0, got %v", v.Type().Field(i).Name, v.Field(i).Float())
		}
	}
}

func TestDSDataAllFieldsSet(t *testing.T) {
	d := DSData{
		Threads:                    1,
		Readwaiters:                2,
		Opsinitiated:               3,
		Opscompleted:               4,
		Dtablesize:                 5,
		Anonymousbinds:             6,
		Unauthbinds:                7,
		Simpleauthbinds:            8,
		Strongauthbinds:            9,
		Bindsecurityerrors:         10,
		Inops:                      11,
		Readops:                    12,
		Compareops:                 13,
		Addentryops:                14,
		Removeentryops:             15,
		Modifyentryops:             16,
		Modifyrdnops:               17,
		Searchops:                  18,
		Onelevelsearchops:          19,
		Wholesubtreesearchops:      20,
		Referrals:                  21,
		Securityerrors:             22,
		Errors:                     23,
		Connections:                24,
		Connectionseq:              25,
		Connectionsinmaxthreads:    26,
		Connectionsmaxthreadscount: 27,
		Bytesrecv:                  28,
		Bytessent:                  29,
		Entriesreturned:            30,
		Referralsreturned:          31,
		Cacheentries:               32,
		Cachehits:                  33,
	}

	v := reflect.ValueOf(d)
	for i := range v.NumField() {
		got := v.Field(i).Float()
		want := float64(i + 1)
		if got != want {
			t.Errorf("%s = %v, want %v", v.Type().Field(i).Name, got, want)
		}
	}
}
