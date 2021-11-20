//+build wireinject

package di

import (
	"github.com/dhis2-sre/go-rate-limiter/pgk/config"
	"github.com/dhis2-sre/go-rate-limiter/pgk/handler"
	"github.com/dhis2-sre/go-rate-limiter/pgk/rule"
	"github.com/google/wire"
	"log"
	"net/http"
)

type Application struct {
	Config  *config.Config
	Handler http.HandlerFunc
}

func ProvideApplication(c *config.Config, h http.HandlerFunc) Application {
	return Application{c, h}
}

func GetApplication() Application {
	wire.Build(
		ProvideApplication,
		provideConfigWithoutError,
		rule.ProvideRouter,
		handler.ProvideHandler,
	)
	return Application{}
}

func provideConfigWithoutError() *config.Config {
	c, err := config.ProvideConfig()
	if err != nil {
		log.Fatal(err)
	}
	return c
}
