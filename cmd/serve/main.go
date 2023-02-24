package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dhis2-sre/go-api-gateway/internal/gateway"
	"github.com/dhis2-sre/go-api-gateway/internal/health"
)

func main() {
	config, err := gateway.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	rules, err := gateway.NewRules(config)
	if err != nil {
		log.Fatal(err)
	}

	router := gateway.NewRouter(rules)

	auth := gateway.NewJwtAuth(config)

	gatewayHandler := gateway.NewHandler(config, router, auth)
	http.HandleFunc("/", gatewayHandler)

	http.HandleFunc("/gateway/health", health.Handler)

	printRules(router.Rules)

	port := config.ServerPort
	server := &http.Server{
		Addr:              ":" + port,
		ReadHeaderTimeout: 3 * time.Second,
	}
	log.Println("Listening on port: " + port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func printRules(rules gateway.Rules) {
	type SetValue struct{}
	ruleSet := map[*gateway.Rule]SetValue{}

	rules.Walk(func(i interface{}) bool {
		rules := i.([]*gateway.Rule)
		for _, rule := range rules {
			ruleSet[rule] = SetValue{}
		}
		return false
	})

	log.Printf("Rules %d (tree: %d)", len(ruleSet), rules.Len())
	for rule := range ruleSet {
		printRule(rule)
	}
}

func printRule(rule *gateway.Rule) {
	method := rule.Method
	if method == "" {
		method = "*"
	}

	logMessage := method
	if rule.Hostname != "" {
		logMessage += fmt.Sprintf(" %s%s -> %s", rule.Hostname, rule.PathPrefix, rule.Backend)
	} else {
		logMessage += fmt.Sprintf(" %s -> %s", rule.PathPrefix, rule.Backend)
	}

	if rule.RequestPerSecond != 0 {
		logMessage += fmt.Sprintf(" - limit(%.2f, %d)", rule.RequestPerSecond, rule.Burst)
	}

	log.Println(logMessage)
}
