tag ?= latest
clean-cmd = docker compose down --remove-orphans --volumes

di:
	wire gen ./pgk/di

binary:
	go build -o go-rate-limit -ldflags "-s -w" ./cmd/serve

docker-image:
	IMAGE_TAG=$(tag) docker-compose build prod && IMAGE_TAG=$(tag) docker-compose push prod

dev:
	docker compose up --build dev backend

test: clean
	docker compose run --no-deps test
	$(clean-cmd)

clean:
	$(clean-cmd)
	go clean

.PHONY: binary docker-image dev
