package yaml

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SSLMode  string `yaml:"sslmode"`
}

type RouteConfig struct {
	DB                 DBConfig          `yaml:"db"`
	UserRoutes         map[string]string `yaml:"user_routes"`
	SubscriptionRoutes map[string]string `yaml:"subscription_routes"`
	PlanRoutes         map[string]string `yaml:"plan_routes"`
	OperatorRoutes     map[string]string `yaml:"operator_routes"`
}
