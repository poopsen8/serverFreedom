package yaml

type DBConfig struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Name       string `yaml:"name"`
	SSLMode    string `yaml:"sslmode"`
	Pathconfig string `yaml:"pathconfig"`
}

type RouteConfig struct {
	DB DBConfig `yaml:"db"`
}
