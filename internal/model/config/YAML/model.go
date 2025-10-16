package yaml

type Config struct {
	Yoomoney *Yoomoney      `yaml:"yoomoney" env-required:"true"`
	Database *DBConfig      `yaml:"db"`
	PathSend *Path          `yaml:"path"`
	Link     ConnectionLink `yaml:"connection_link"`
}

type DBConfig struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Name       string `yaml:"name"`
	SSLMode    string `yaml:"sslmode"`
	Pathconfig string `yaml:"pathconfig"`
}

type Yoomoney struct {
	BaseURL  string `yaml:"base_url" env-required:"true"`
	BasePath string `yaml:"base_path" env-required:"true"`

	Receiver *YoomoneyReceiver `yaml:"receiver" env-required:"true"`
}

type YoomoneyReceiver struct {
	Account     string `yaml:"account" env-required:"true"`
	NotifSecret string `yaml:"notification_secret" env-required:"true"`
}

type Path struct {
	Base_url string `yaml:"base_url"`
}
type ConnectionLink struct {
	Left_part  string `yaml:"left_part"`
	Right_part string `yaml:"right_part"`
}
