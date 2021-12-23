package main

import (
	"fmt"
	"github.com/dhis2-sre/go-api-gateway/internal/gateway"
	"github.com/dhis2-sre/go-api-gateway/internal/health"
	"log"
	"net/http"
)

func main() {
	config, err := gateway.ProvideConfig()
	if err != nil {
		log.Fatal(err)
	}

	rules, err := gateway.ProvideRules(config)
	if err != nil {
		log.Fatal(err)
	}

	router := gateway.ProvideRouter(rules)

	gatewayHandler := gateway.ProvideHandler(config, router)
	http.HandleFunc("/", gatewayHandler)

	http.HandleFunc("/gateway/health", health.Handler)

	printRules(router)

	port := config.ServerPort
	log.Println("Listening on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func printRules(router *gateway.Router) {
	type SetValue struct{}
	ruleSet := map[*gateway.Rule]SetValue{}

	router.Rules.Root().Walk(func(k []byte, i interface{}) bool {
		rules := i.([]*gateway.Rule)
		for _, rule := range rules {
			ruleSet[rule] = SetValue{}
		}
		return false
	})

	log.Printf("Rules %d (tree: %d)", len(ruleSet), router.Rules.Len())
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
