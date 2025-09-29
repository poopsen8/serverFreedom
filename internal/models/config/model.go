package config

type Config struct {
	Inbounds  []Inbound  `json:"inbounds"`
	Outbounds []Outbound `json:"outbounds"`
	Routing   Routing    `json:"routing"`
}

type Inbound struct {
	Listen         string         `json:"listen"`
	Port           int            `json:"port"`
	Protocol       string         `json:"protocol"`
	Settings       Settings       `json:"settings"`
	StreamSettings StreamSettings `json:"streamSettings"`
	Sniffing       Sniffing       `json:"sniffing"`
}

type Settings struct {
	Clients    []Client `json:"clients"`
	Decryption string   `json:"decryption"`
}

type Client struct {
	ID   string `json:"id"`
	Flow string `json:"flow"`
}

type StreamSettings struct {
	Network         string          `json:"network"`
	Security        string          `json:"security"`
	RealitySettings RealitySettings `json:"realitySettings"`
}

type RealitySettings struct {
	Show         bool     `json:"show"`
	Dest         string   `json:"dest"`
	Xver         int      `json:"xver"`
	ServerNames  []string `json:"serverNames"`
	PrivateKey   string   `json:"privateKey"`
	MinClientVer string   `json:"minClientVer"`
	MaxClientVer string   `json:"maxClientVer"`
	MaxTimeDiff  int      `json:"maxTimeDiff"`
	ShortIds     []string `json:"shortIds"`
}

type Sniffing struct {
	Enabled      bool     `json:"enabled"`
	DestOverride []string `json:"destOverride"`
	MetadataOnly bool     `json:"metadataOnly"`
	RouteOnly    bool     `json:"routeOnly"`
}

type Outbound struct {
	Protocol       string                 `json:"protocol"`
	Tag            string                 `json:"tag"`
	Settings       OutboundSettings       `json:"settings"`
	StreamSettings OutboundStreamSettings `json:"streamSettings"`
}

type OutboundSettings struct {
	DomainStrategy string `json:"domainStrategy"`
	Redirect       string `json:"redirect"`
}

type OutboundStreamSettings struct {
	Sockopt Sockopt `json:"sockopt"`
}

type Sockopt struct {
	TcpFastOpen bool   `json:"tcpFastOpen"`
	Tproxy      string `json:"tproxy"`
}

type Routing struct {
	DomainStrategy string `json:"domainStrategy"`
	Rules          []Rule `json:"rules"`
}

type Rule struct {
	Type        string   `json:"type"`
	IP          []string `json:"ip"`
	OutboundTag string   `json:"outboundTag"`
}
