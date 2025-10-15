package json

type SingBoxConfig struct {
	Log          Log          `json:"log"`
	DNS          DNS          `json:"dns"`
	Inbounds     []Inbound    `json:"inbounds"`
	Outbounds    []Outbound   `json:"outbounds"`
	Route        Route        `json:"route"`
	Experimental Experimental `json:"experimental"`
}

type Log struct {
	Level     string `json:"level"`
	Timestamp bool   `json:"timestamp"`
}

type DNS struct {
	Servers  []DNSServer `json:"servers"`
	Strategy string      `json:"strategy"`
}

type DNSServer struct {
	Address string `json:"address"`
	Detour  string `json:"detour"`
}

type Inbound struct {
	Type                     string `json:"type"`
	Listen                   string `json:"listen"`
	ListenPort               int    `json:"listen_port"`
	Network                  string `json:"network,omitempty"`
	Sniff                    bool   `json:"sniff,omitempty"`
	SniffOverrideDestination bool   `json:"sniff_override_destination,omitempty"`
	DomainStrategy           string `json:"domain_strategy,omitempty"`
	OverrideAddress          string `json:"override_address,omitempty"`
	OverridePort             int    `json:"override_port,omitempty"`

	// Trojan-specific
	Users []User `json:"users,omitempty"`
	TLS   *TLS   `json:"tls,omitempty"`
}

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type TLS struct {
	Enabled    bool     `json:"enabled"`
	ServerName string   `json:"server_name"`
	ALPN       []string `json:"alpn"`
	Reality    *Reality `json:"reality"`
}

type Reality struct {
	Enabled           bool      `json:"enabled"`
	Handshake         Handshake `json:"handshake"`
	PrivateKey        string    `json:"private_key"`
	ShortID           []string  `json:"short_id"`
	MaxTimeDifference string    `json:"max_time_difference"`
}

type Handshake struct {
	Server     string `json:"server"`
	ServerPort int    `json:"server_port"`
}

type Outbound struct {
	Type string `json:"type"`
	Tag  string `json:"tag"`
}

type Route struct {
	Final   string    `json:"final"`
	RuleSet []RuleSet `json:"rule_set"`
	Rules   []Rule    `json:"rules"`
}

type RuleSet struct {
	Tag            string `json:"tag"`
	Type           string `json:"type"`
	Format         string `json:"format"`
	URL            string `json:"url"`
	DownloadDetour string `json:"download_detour"`
}

type Rule struct {
	RuleSet  []string `json:"rule_set,omitempty"`
	Network  string   `json:"network,omitempty"`
	Port     []int    `json:"port,omitempty"`
	Outbound string   `json:"outbound"`
}

type Experimental struct {
	CacheFile CacheFile `json:"cache_file"`
}

type CacheFile struct {
	Enabled bool `json:"enabled"`
}
