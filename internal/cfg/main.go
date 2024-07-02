package cfg

type Config struct {
	Monitors []MonitorConfig `yaml:"monitors"`
}

type MonitorConfig struct {
	Name string `yaml:"name"`
	Url string `yaml:"url"`
	Type string `yaml:"type"`
	Heartbeat string `yaml:"heartbeat"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Expects []MonitorExpects `yaml:"expects"`
	Failures []string `yaml:"failures"`
}

type MonitorExpects struct {
	Field string `yaml:"field"`
	Op string `yaml:"op"`
	Value string `yaml:"value"`
}
