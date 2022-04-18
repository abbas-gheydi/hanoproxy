package model

const (
	Postgres      string = "postgres"
	Sentinel      string = "sentinel"
	None          string = ""
	HTTP          string = "http"
	TCP           string = "tcp"
	ActivePassive string = "active-passive"
	RoundRobin    string = "roundrobin"
)

type IpAddr struct {
	Addr      string
	Port      string
	IsHealthy bool
	UserName  string
	Password  string
	Url       string
}

type Record struct {
	Name        string
	Ip          []IpAddr
	Sentinels   []IpAddr
	LastHint    int
	ServiceType string
	Options     Option
}

type Option struct {
	MasterOnly                bool
	SlaveOnly                 bool
	CheckForHealth            bool
	DbName                    string
	SentinelMonitorMasterName string
	Expected_Response_Code    int
	HttpMethod                string
	RetryCount                int
	LBmethod                  string
}
