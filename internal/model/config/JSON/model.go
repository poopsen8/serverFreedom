package config

type SingBoxConfig struct {
	Log          LogConfig    `json:"log"`
	DNS          DNSConfig    `json:"dns"`
	Inbounds     []Inbound    `json:"inbounds"`
	Outbounds    []Outbound   `json:"outbounds"`
	Route        RouteConfig  `json:"route"`
	Experimental Experimental `json:"experimental"`
}

type LogConfig struct {
	Level     string `json:"level"`
	Timestamp bool   `json:"timestamp"`
}

type DNSConfig struct {
	Servers  []DNSServer `json:"servers"`
	Strategy string      `json:"strategy"`
}

type DNSServer struct {
	Address string `json:"address"`
	Detour  string `json:"detour"`
}

type Inbound struct {
	Type                     string     `json:"type"`
	Listen                   string     `json:"listen"`
	ListenPort               int        `json:"listen_port"`
	Network                  string     `json:"network,omitempty"`
	OverrideAddress          string     `json:"override_address,omitempty"`
	OverridePort             int        `json:"override_port,omitempty"`
	Sniff                    bool       `json:"sniff,omitempty"`
	SniffOverrideDestination bool       `json:"sniff_override_destination,omitempty"`
	DomainStrategy           string     `json:"domain_strategy,omitempty"`
	Users                    []User     `json:"users,omitempty"`
	TLS                      *TLSConfig `json:"tls,omitempty"`
}

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type TLSConfig struct {
	Enabled    bool           `json:"enabled"`
	ServerName string         `json:"server_name"`
	ALPN       []string       `json:"alpn"`
	Reality    *RealityConfig `json:"reality"`
}

type RealityConfig struct {
	Enabled           bool            `json:"enabled"`
	Handshake         HandshakeConfig `json:"handshake"`
	PrivateKey        string          `json:"private_key"`
	ShortID           []string        `json:"short_id"`
	MaxTimeDifference string          `json:"max_time_difference"`
}

type HandshakeConfig struct {
	Server     string `json:"server"`
	ServerPort int    `json:"server_port"`
}

type Outbound struct {
	Type string `json:"type"`
	Tag  string `json:"tag"`
}

type RouteConfig struct {
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
	CacheFile CacheFileConfig `json:"cache_file"`
}

type CacheFileConfig struct {
	Enabled bool `json:"enabled"`
}
