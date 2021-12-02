package main

import (
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

	router, err := gateway.ProvideRouter(config)
	if err != nil {
		log.Fatal(err)
	}

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

	if router.CatchAllRule != nil {
		ruleSet[router.CatchAllRule] = SetValue{}
	}

	log.Printf("Rules %d (tree: %d)", len(ruleSet), router.Rules.Len())
	for rule, _ := range ruleSet {
		printRule(rule)
	}
}

func printRule(rule *gateway.Rule) {
	method := rule.Method
	if method == "" {
		method = "*"
	}

	if rule.RequestPerSecond != 0 {
		log.Printf("%s %s -> %s - limit(%.2f, %d)", method, rule.PathPrefix, rule.Backend, rule.RequestPerSecond, rule.Burst)
	} else {
		log.Printf("%s %s -> %s", method, rule.PathPrefix, rule.Backend)
	}
}
