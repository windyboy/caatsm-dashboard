package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	ProductionEnv  = "prod"
	DevelopmentEnv = "dev"

	KeyEnv = "GO_ENV"
)

type Config struct {
	Server       ServerConfig     `mapstructure:"server"`
	Nats         NatsConfig       `mapstructure:"nats"`
	Timeouts     TimeoutsConfig   `mapstructure:"timeouts"`
	Subscription SubscriberConfig `mapstructure:"subscription"`
	Hasura       HasuraConfig     `mapstructure:"hasura"`
}

type ServerConfig struct {
	ListenAddr string `mapstructure:"listenAddr"`
}

type NatsConfig struct {
	URL     string `mapstructure:"url"`
	Client  string `mapstructure:"client"`
	Cluster string `mapstructure:"cluster"`
}

type SubscriberConfig struct {
	Topic string `mapstructure:"topic"`
	Queue string `mapstructure:"queue"`
}

type TimeoutsConfig struct {
	Server        time.Duration `mapstructure:"server"`
	ReconnectWait time.Duration `mapstructure:"reconnect_wait"`
	Close         time.Duration `mapstructure:"close"`
	AckWait       time.Duration `mapstructure:"ack_wait"`
}

type HasuraConfig struct {
	Endpoint string `mapstructure:"endpoint"`
	Secret   string `mapstructure:"secret"`
}

// LoadConfig loads the configuration from a file
func LoadConfig() (*Config, error) {
	// log := utils.Logger
	env := os.Getenv(KeyEnv)
	if env == "" {
		env = DevelopmentEnv
	}

	viper.SetConfigType("toml")
	viper.SetConfigName("config." + env)
	viper.AddConfigPath("configs")
	viper.SetEnvPrefix("tele")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file for environment %q: %w", env, err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct for environment %q: %w", env, err)
	}
	if err := config.Validate(); err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *Config) Validate() error {
	var missing []string

	if c.Server.ListenAddr == "" {
		missing = append(missing, "server.listenAddr")
	}

	if c.Nats.URL == "" {
		missing = append(missing, "nats.url")
	}

	if c.Subscription.Topic == "" {
		missing = append(missing, "subscription.topic")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required configuration values: %s", strings.Join(missing, ", "))
	}
	return nil
}
