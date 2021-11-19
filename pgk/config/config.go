package config

import (
	"github.com/spf13/viper"
)

func ProvideConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	var c *Config
	if err := viper.ReadInConfig(); err != nil {
		return &Config{}, err
	}

	err := viper.Unmarshal(&c)
	if err != nil {
		return &Config{}, err
	}

	return c, nil
}

type Config struct {
	ServerPort     string
	Authentication Authentication
	Rules          []Rule
}

type Authentication struct {
	Jwt Jwt
}

type Jwt struct {
	PublicKey string
}

type Rule struct {
	Method           string
	PathPattern      string
	Backend          string
	Authentication   string
	RequestPerSecond float64
	Burst            int
}
