package boot

import (
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var envOnce sync.Once

var config Config

// Config holds all the application configuration parameters. With the help of
// 'envconfig' module we load them from ENV, also doing some basic validation.
type Config struct {
	AppName                 string `split_words:"true" default:"unknown"`
	AppBanner               bool   `split_words:"true" default:"false"`
	AppInitOutboxDispatcher bool   `split_words:"true" default:"false"`

	LogLevel    int  `split_words:"true" default:"1"`
	LogBeautify bool `split_words:"true" default:"false"`

	DbHost     string `split_words:"true" required:"true"`
	DbPort     string `split_words:"true" required:"true"`
	DbName     string `split_words:"true" required:"true"`
	DbUser     string `split_words:"true" required:"true"`
	DbPassword string `split_words:"true" required:"true"`

	GinMode string `envconfig:"GIN_MODE" split_words:"true" default:"release"`

	KafkaBootstrapServers string `split_words:"true"`
	KafkaSchemaRegistry   string `split_words:"true" required:"true"`
}

// LoadConfig bootstraps the application configuration (setting the ones not
// already present in the environment reading from .env, and loading them into
// the Config struct for better accessibility from the adapter layers).
func LoadConfig() {
	envOnce.Do(func() {
		loadEnvs()
		mapEnvsToConfig()
	})
}

// GetConfig provides the a reference to the configuration struct.
func GetConfig() *Config {
	return &config
}

// loadEnvs uses 'godotenv' to preload the application configuration into env
// variables. It only sets the ones already not set, so no overriding is done.
func loadEnvs() {
	godotenv.Load()
}

// mapEnvsToConfig uses 'envconfig' to map env variables to a configuration struct for better
// access from the adapter layers (it also panics when required parameters are not present).
func mapEnvsToConfig() {
	config = Config{}
	envconfig.MustProcess("f4allgo", &config)
}
