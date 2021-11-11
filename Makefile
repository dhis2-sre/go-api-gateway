tag ?= latest

binary:
	go build -o go-rate-limit -ldflags "-s -w" .

docker-image:
	IMAGE_TAG=$(tag) docker-compose build prod && IMAGE_TAG=$(tag) docker-compose push prod

dev:
	docker compose up --build dev

run:
	go run .

.PHONY: binary docker-image
