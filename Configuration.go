package main

type Configuration struct {
	ServerPort string
	Backends   []string
	Rules      []Rule
}
