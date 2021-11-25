tag ?= latest
version ?= $(shell yq e '.version' helm/Chart.yaml)
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

helm-chart:
	@helm package helm/chart

publish-helm:
	@curl --user "$(CHART_AUTH_USER):$(CHART_AUTH_PASS)" \
        -F "chart=@api-gateway-$(version).tgz" \
        -F "prov=@api-gateway-$(version).tgz.prov" \
        https://helm-charts.fitfit.dk/api/charts

.PHONY: binary docker-image push-docker-image dev test dev-test helm-package publish-helm
