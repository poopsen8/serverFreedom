package yaml

type RouteConfig struct {
	UR UserRoutes         `yaml:"user_routes" json:"user_routes"`
	SR SubscriptionRoutes `yaml:"subscription_routes" json:"subscription_routes"`
	PR PlanRoutes         `yaml:"plan_routes" json:"plan_routes"`
	OP OperatorRoutes     `yaml:"operator_routes" json:"operator_routes"`
}

type UserRoutes struct {
	Register string `yaml:"register" json:"register"`
	Update   string `yaml:"update" json:"update"`
	Get      string `yaml:"get" json:"get"`
}

type SubscriptionRoutes struct {
	AddSubscription string `yaml:"add_subscription"`
	UpdateKey       string `yaml:"update_key"`
	Get             string `yaml:"get"`
}

type PlanRoutes struct {
	GetAll string `yaml:"get_all"`
	Get    string `yaml:"get"`
}

type OperatorRoutes struct {
	GetAll string `yaml:"get_all"`
	Get    string `yaml:"get"`
}
