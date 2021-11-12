tag ?= latest

di:
	wire gen ./pgk/di

binary:
	go build -o go-rate-limit -ldflags "-s -w" ./cmd/serve

docker-image:
	IMAGE_TAG=$(tag) docker-compose build prod && IMAGE_TAG=$(tag) docker-compose push prod

dev:
	docker compose up --build dev backend

.PHONY: binary docker-image dev
