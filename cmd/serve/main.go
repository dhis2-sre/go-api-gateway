package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dhis2-sre/go-rate-limite/pgk/config"
	"github.com/dhis2-sre/go-rate-limite/pgk/handler"
	"github.com/dhis2-sre/go-rate-limite/pgk/proxy"
	"github.com/dhis2-sre/go-rate-limite/pgk/rule"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed due to: %s", err)
		os.Exit(1)
	}
}

func run() error {
	c, err := config.ProvideConfig()
	if err != nil {
		return err
	}

	rules := rule.ProvideRules(c)
	s := &http.Server{
		Addr:    ":" + c.ServerPort,
		Handler: handler.RateLimit(rules)(proxy.TransparentProxy(c.Backend)),
	}
	log.Println("Listening on port: " + c.ServerPort)
	return s.ListenAndServe()
}
