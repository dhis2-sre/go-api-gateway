package main

import (
	"github.com/dhis2-sre/go-rate-limiter/pgk/di"
	"log"
	"net/http"
)

func main() {
	app := di.GetApplication()

	port := app.Config.ServerPort
	handler := http.HandlerFunc(app.Handler.RateLimitingProxyHandler)

	log.Println("Listening on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
