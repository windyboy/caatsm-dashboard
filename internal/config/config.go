package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	ProdcutionEnv  = "prod"
	DevelopmentEnv = "dev"

	KeyEnv = "GO_ENV"
)

type Config struct {
	Server       ServerConfig
	Nats         NatsConfig
	Timeouts     TimeoutsConfig
	Subscription SubscriberConfig
}

type ServerConfig struct {
	ListenAddr string
}

type NatsConfig struct {
	URL string
}

type SubscriberConfig struct {
	Topic string `mapstructure:"topic"`
}

type TimeoutsConfig struct {
	Server        time.Duration `mapstructure:"server"`
	ReconnectWait time.Duration `mapstructure:"reconnect_wait"`
	Close         time.Duration `mapstructure:"close"`
	AckWait       time.Duration `mapstructure:"ack_wait"`
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

	if err := viper.ReadInConfig(); err != nil {
		errMsg := fmt.Sprintf("error reading config file for environment '%s': %v", env, err)
		// log.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		errMsg := fmt.Sprintf("unable to decode config into struct for environment '%s': %v", env, err)
		// log.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	return &config, nil
}
