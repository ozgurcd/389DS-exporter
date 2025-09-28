package obj

// DSData stores metrics from 389DS
type DSData struct {
	Threads                    float64
	Readwaiters                float64
	Opsinitiated               float64
	Opscompleted               float64
	Dtablesize                 float64
	Anonymousbinds             float64
	Unauthbinds                float64
	Simpleauthbinds            float64
	Strongauthbinds            float64
	Bindsecurityerrors         float64
	Inops                      float64
	Readops                    float64
	Compareops                 float64
	Addentryops                float64
	Removeentryops             float64
	Modifyentryops             float64
	Modifyrdnops               float64
	Searchops                  float64
	Onelevelsearchops          float64
	Wholesubtreesearchops      float64
	Referrals                  float64
	Securityerrors             float64
	Errors                     float64
	Connections                float64
	Connectionseq              float64
	Connectionsinmaxthreads    float64
	Connectionsmaxthreadscount float64
	Bytesrecv                  float64
	Bytessent                  float64
	Entriesreturned            float64
	Referralsreturned          float64
	Cacheentries               float64
	Cachehits                  float64
}
