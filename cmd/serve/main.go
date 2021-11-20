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

	port := config.ServerPort
	log.Println("Listening on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
