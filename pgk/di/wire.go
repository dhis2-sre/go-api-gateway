//+build wireinject

package di

import (
	"github.com/dhis2-sre/go-rate-limiter/pgk/config"
	"github.com/dhis2-sre/go-rate-limiter/pgk/handler"
	"github.com/dhis2-sre/go-rate-limiter/pgk/proxy"
	"github.com/dhis2-sre/go-rate-limiter/pgk/rule"
	"github.com/google/wire"
	"log"
)

type Application struct {
	Config  *config.Config
	Handler handler.Handler
}

func ProvideApplication(c *config.Config, h handler.Handler) Application {
	return Application{c, h}
}

func GetApplication() Application {
	wire.Build(
		ProvideApplication,
		provideConfigWithoutError,
		proxy.ProvideProxy,
		rule.ProvideRules,
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
