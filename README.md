# Introduction
Lightweight HTTP API Gateway.

The gateway will perform proxying, rate limiting and token validation based on the rules defined in the configuration [file](config.yml).

# Architecture
## Radix
## Header matching

# Development
## Launch in development mode
Any changes made to the go or yml files will result in the application being recompiled and relaunched
```sh
make dev
```

## Launch in unit test development mode
Just like `make dev` tests will be rerun if any changes are made to either go or yml files 
```sh
make dev-test
```

## Kubernetes
### Helm

### Helmfile

### Skaffold
```sh
skaffold dev
```

# HTTP Request
Perform a http request using HTTPie
```sh
http :8080
```

# Trigger rate limiting
```sh
seq 10 | xargs -P 4 -I '{}' http post :8080/health
```

# Configuration
## Minimal

## Rate limiting

## Request Matching
### Path Prefix
### HTTP Method
### HTTP Headers

## Environment variables



# Credits
* https://github.com/didip/tollbooth

# TODO
* Comment exported types and functions
* Write tests
* CICD
