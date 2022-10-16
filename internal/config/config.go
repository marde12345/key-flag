package config

type Config struct {
	Resources Resources `yaml:"resources"`
}

type Resources struct {
	Redis  Redis  `yaml:"redis"`
	Consul Consul `yaml:"consul"`
	DB     DB     `yaml:"db"`
	Slack  Slack  `yaml:"slack"`
}

type Slack struct {
	Webhook string `yaml:"webhook"`
}

type Redis struct {
	Address string `yaml:"address"`
}

type Consul struct {
	Address string `yaml:"address"`
}

type DB struct {
	MasterAddress   string `yaml:"masterAddress"`
	FollowerAddress string `yaml:"followerAddress"`
}

var KVMiddlewareConfig Config
