package main

import (
	"fmt"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
)

var backendCount = 0

func loadBalance(backends []string) string {
	backend := backends[backendCount]
	backendCount++

	if backendCount >= len(backends) {
		backendCount = 0
	}

	return backend
}

func main() {
	configuration := readConfiguration()

	// TODO: This is only defined here because "configuration" is needed
	handleRequest := func(res http.ResponseWriter, req *http.Request) {
		backend := loadBalance(configuration.Backends)
		log.Printf("%s %s -> %s", req.Method, req.URL.Path, backend)

		backendUrl, _ := url.Parse(backend)
		proxy := httputil.NewSingleHostReverseProxy(backendUrl)
		proxy.ServeHTTP(res, req)
	}

	rules := loadRules(configuration)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if match, lmt := rules.match(r); match {
			handler := tollbooth.LimitFuncHandler(lmt, handleRequest)
			handler.ServeHTTP(w, r)
		} else {
			handleRequest(w, r)
		}
	})

	log.Println("Listening on port: " + configuration.ServerPort)
	log.Fatal(http.ListenAndServe(":"+configuration.ServerPort, nil))
}

func loadRules(configuration Configuration) *Rules {
	for i, rule := range configuration.Rules {
		lmt := tollbooth.NewLimiter(rule.RequestPerSecond, nil)
		lmt.SetMethods([]string{rule.Method})
		lmt.SetBurst(rule.Burst)
		rule.Lmt = lmt

		configuration.Rules[i] = rule
	}

	rules := &Rules{
		Rules: configuration.Rules,
	}
	return rules
}

func readConfiguration() Configuration {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	var configuration Configuration
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	return configuration
}

type Rules struct {
	Rules []Rule
}

func (r Rules) match(req *http.Request) (bool, *limiter.Limiter) {
	for _, rule := range r.Rules {
		if rule.pathMatch(req.URL.Path) {
			return true, rule.Lmt
		}
	}
	return false, &limiter.Limiter{}
}

type Rule struct {
	Method           string
	PathPattern      string
	RequestPerSecond float64
	Burst            int
	Lmt              *limiter.Limiter
}

func (r *Rule) pathMatch(path string) bool {
	match, err := regexp.MatchString(r.PathPattern, path)
	if err != nil {
		return false
	}
	return match
}
