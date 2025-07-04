package config

type EnvConfig struct {
	EnvironmentEnv EnvironmentConfig `mapstructure:"ENVIRONMENT"`
	ServerEnv      ServerConfig      `mapstructure:"SERVER"`
	MongoDBEnv     MongoConfig       `mapstructure:"MONGODB"`
	LoggerEnv      LoggerConfig      `mapstructure:"LOGGER"`
	GenieacsEnv    GeniacsConfig     `mapstructure:"GENIEACS"`
	AuthEnv        Auth              `mapstructure:"AUTH"`
}

type EnvironmentConfig struct {
	Env     string `mapstructure:"ENV"`
	App     string `mapstructure:"APP"`
	Version string `mapstructure:"VERSION"`
	Host    string `mapstructure:"HOST"`
}

type GeniacsConfig struct {
	NBI_URL string `mapstructure:"NBI_URL"`
}

type Auth struct {
	APIKey string `mapstructure:"APIKEY"`
}

type LoggerConfig struct {
	LogLevel string `mapstructure:"LOGLEVEL"`
}

type ServerConfig struct {
	Port        string `mapstructure:"PORT"` // Consistent with YAML string
	MetricsPort string `mapstructure:"METRICSPORT"`
}

type MongoConfig struct {
	URI        string `mapstructure:"URI"`
	Host       string `mapstructure:"HOST"`
	Port       string `mapstructure:"PORT"` // Changed to string for consistency
	TLSEnabled bool   `mapstructure:"TLSENABLED"`
	Database   string `mapstructure:"DATABASE"`
	Username   string `mapstructure:"USERNAME"`
	Password   string `mapstructure:"PASSWORD"`
}
