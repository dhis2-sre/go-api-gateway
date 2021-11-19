package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dhis2-sre/go-rate-limite/pgk/config"
	"github.com/dhis2-sre/go-rate-limite/pgk/proxy"
	"github.com/dhis2-sre/go-rate-limite/pgk/rate"
	"github.com/dhis2-sre/go-rate-limite/pgk/rule"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed due to: %s", err)
		os.Exit(1)
	}
}

// TODO: for testability you can pass in the io.Writer for logging, and any
// os.Args if needed
func run() error {
	c, err := config.ProvideConfig()
	if err != nil {
		return err
	}

	rules := rule.NewRules(c)
	// TODO: define any timeouts if you want
	s := &http.Server{
		Addr:    ":" + c.ServerPort,
		Handler: rate.Limit(rules)(proxy.Transparently(c.Backend)),
	}
	log.Println("Listening on port: " + c.ServerPort)
	return s.ListenAndServe()
}
