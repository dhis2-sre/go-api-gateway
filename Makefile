tag ?= latest
clean-cmd = docker compose down --remove-orphans --volumes

init:
	pip install pre-commit
	pre-commit install --install-hooks --overwrite

	go install github.com/direnv/direnv@latest
	direnv version

	go install golang.org/x/tools/cmd/goimports@latest

	go install github.com/securego/gosec/v2/cmd/gosec@latest
	gosec --version

check:
	pre-commit run --all-files --show-diff-on-failure

docker-image:
	IMAGE_TAG=$(tag) docker compose build prod

push-docker-image:
	IMAGE_TAG=$(tag) docker compose push prod

smoke-test:
	IMAGE_TAG=$(tag) docker compose up -d prod

dev:
	docker compose up --build dev backend0 backend1 jwks

test: clean
	docker compose up -d backend0 backend1 jwks
	docker compose run --no-deps test
	$(clean-cmd)

dev-test: clean
	docker compose up -d backend0 backend1 jwks
	docker compose run --no-deps dev-test
	$(clean-cmd)

clean:
	$(clean-cmd)
	go clean

.PHONY: init check docker-image push-docker-image dev test dev-test
