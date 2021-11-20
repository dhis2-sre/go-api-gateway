package main

import (
	"github.com/dhis2-sre/go-rate-limiter/internal/gateway"
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

	handler := gateway.ProvideHandler(config, router)

	printRules(router)

	port := config.ServerPort
	log.Println("Listening on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func printRules(router *gateway.Router) {
	log.Println("Rules:")
	for _, rule := range router.Rules {
		method := rule.Method
		if method == "" {
			method = "*"
		}
		log.Printf("%s %s -> %s - limit(%.2f, %d)", method, rule.PathPattern, rule.Backend, rule.RequestPerSecond, rule.Burst)
	}
}
