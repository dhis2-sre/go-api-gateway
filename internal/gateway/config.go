package gateway

import (
	"net/http"
	"time"

	"github.com/spf13/viper"
)

func NewConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		return &Config{}, err
	}

	bindMap := map[string]string{
		"serverport":                   "APIG_SERVER_PORT",
		"basepath":                     "APIG_BASE_PATH",
		"maxmultipartsize":             "APIG_MAX_MULTIPART_SIZE",
		"defaultbackend":               "APIG_DEFAULT_BACKEND",
		"authentication.jwt.publickey": "APIG_PUBLIC_KEY",
	}

	for k, v := range bindMap {
		err := viper.BindEnv(k, v)
		if err != nil {
			return &Config{}, err
		}
	}

	viper.SetDefault("maxmultipartsize", 2)

	var c *Config
	err := viper.Unmarshal(&c)
	if err != nil {
		return &Config{}, err
	}

	return c, nil
}

type Config struct {
	ServerPort       string
	BasePath         string
	MaxMultipartSize int64
	DefaultBackend   string
	Authentication   Authentication
	Backends         []Backend
	Rules            []ConfigRule
}

type Authentication struct {
	Jwt  Jwt
	Jwks Jwks
}

type Jwt struct {
	PublicKey string
}

type Jwks struct {
	Host                   string
	Index                  int
	MinimumRefreshInterval time.Duration
}

type Backend struct {
	Name string
	Url  string
}

type ConfigRule struct {
	Method           string
	PathPrefix       string
	Hostname         string
	PathReplace      PathReplace
	Block            bool
	Backend          string
	Authentication   string
	RequestPerSecond float64
	Burst            int
	Headers          http.Header
}

type PathReplace struct {
	Target      string
	Replacement string
}
