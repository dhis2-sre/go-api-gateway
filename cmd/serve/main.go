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
	log.Printf("Rules (%d)", router.Rules.Len())
	router.Rules.Root().Walk(func(k []byte, i interface{}) bool {
		rule := i.(*gateway.Rule)
		method := rule.Method
		if method == "" {
			method = "*"
		}

		if rule.RequestPerSecond != 0 {
			log.Printf("%s %s -> %s - limit(%.2f, %d)", method, rule.PathPrefix, rule.Backend, rule.RequestPerSecond, rule.Burst)
		} else {
			log.Printf("%s %s -> %s", method, rule.PathPrefix, rule.Backend)
		}
		return false
	})
}
