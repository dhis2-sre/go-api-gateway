tag ?= latest
clean-cmd = docker compose down --remove-orphans --volumes

binary:
	go build -o go-api-gateway -ldflags "-s -w" ./cmd/serve

docker-image:
	IMAGE_TAG=$(tag) docker compose build prod

push-docker-image:
	IMAGE_TAG=$(tag) docker compose push prod

dev:
	docker compose up --build dev backend0 backend1

test: clean
	docker compose up -d backend0 backend1
	docker compose run --no-deps test
	$(clean-cmd)

dev-test: clean
	docker compose up -d backend0 backend1
	docker compose run --no-deps dev-test
	$(clean-cmd)

clean:
	$(clean-cmd)
	go clean

.PHONY: binary docker-image dev
