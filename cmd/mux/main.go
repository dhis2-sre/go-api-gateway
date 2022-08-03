package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"time"

	"github.com/dhis2-sre/go-api-gateway/internal/gateway"
	"github.com/dhis2-sre/go-api-gateway/internal/health"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	config, err := gateway.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	rules, err := gateway.NewRuleList(config)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/gateway/health", health.Handler)

	addRules(r, config, rules)

	printRules(rules)

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	srv := &http.Server{
		Handler:      loggedRouter,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv.ListenAndServe()
}

func max(i, j string) string {
	if len(i) > len(j) {
		return i
	}
	return j
}

func addRules(r *mux.Router, c *gateway.Config, rules []gateway.ConfigRule) {
	auth := gateway.NewJwtAuth(c)

	// TODO: Sort rules by len(pathPrefix, path) where path is s/{*}//
	sort.SliceStable(rules, func(i, j int) bool {
		ri := rules[i]
		rj := rules[j]
		re := regexp.MustCompile(`{[^}]*}`)
		iPath := re.ReplaceAllString(ri.Path, "")
		jPath := re.ReplaceAllString(rj.Path, "")
		si := max(iPath, ri.PathPrefix)
		sj := max(jPath, rj.PathPrefix)
		return len(si) > len(sj)
	})

	for _, rule := range rules {
		methods := getMethods(rule)
		route := r.Methods(methods...)

		if rule.Path != "" {
			route.Path(rule.Path)
		}

		if rule.PathPrefix != "" {
			route.PathPrefix(rule.PathPrefix)
		}

		if rule.Hostname != "" {
			route.Host(rule.Hostname)
		}

		for key, values := range rule.Headers {
			for _, value := range values {
				route.Headers(key, value)
			}
		}

		if rule.Block {
			route.HandlerFunc(gateway.NewBlockingProxy)
		}

		handler2, err := gateway.NewHandler2(rule, auth)
		if err != nil {
			log.Fatal(err)
		}
		route.Handler(handler2)
	}
}

func getMethods(rule gateway.ConfigRule) []string {
	if rule.Method != "" {
		return []string{rule.Method}
	}
	return []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace}
}

func printRules(rules []gateway.ConfigRule) {
	log.Printf("Rules %d", len(rules))
	for _, rule := range rules {
		method := rule.Method
		if method == "" {
			method = "*"
		}

		logMessage := method
		if rule.Hostname != "" {
			logMessage += fmt.Sprintf(" %s%s -> %s", rule.Hostname, max(rule.Path, rule.PathPrefix), rule.Backend)
		} else {
			logMessage += fmt.Sprintf(" %s -> %s", max(rule.Path, rule.PathPrefix), rule.Backend)
		}

		if rule.RequestPerSecond != 0 {
			logMessage += fmt.Sprintf(" - limit(%.2f, %d)", rule.RequestPerSecond, rule.Burst)
		}

		log.Println(logMessage)
	}
}
