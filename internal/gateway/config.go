package gateway

import (
	"github.com/spf13/viper"
	"net/http"
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

	bindMap := map[string]string{
		"serverport":                   "APIG_SERVER_PORT",
		"basepath":                     "APIG_BASE_PATH",
		"defaultbackend":               "APIG_DEFAULT_BACKEND",
		"authentication.jwt.publickey": "APIG_PUBLIC_KEY",
	}

	for k, v := range bindMap {
		err := viper.BindEnv(k, v)
		if err != nil {
			return &Config{}, err
		}
	}

	err := viper.Unmarshal(&c)
	if err != nil {
		return &Config{}, err
	}

	return c, nil
}

type Config struct {
	ServerPort     string
	BasePath       string
	DefaultBackend string
	Authentication Authentication
	Backends       []Backend
	Rules          []ConfigRule
}

type Authentication struct {
	Jwt Jwt
}

type Jwt struct {
	PublicKey string
}

type Backend struct {
	Name string
	Url  string
}

type ConfigRule struct {
	Method           string
	PathPrefix       string
	Block            bool
	Backend          string
	Authentication   string
	RequestPerSecond float64
	Burst            int
	Headers          http.Header
}
