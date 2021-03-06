# Introduction

Lightweight HTTP API Gateway.

The gateway will perform proxying, rate limiting and token validation based on the rules defined in the
configuration [file](config.yml).

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [Dependencies (for development)](#dependencies-for-development)
- [Quick start](#quick-start)
- [Development](#development)
  - [Unit Test Development](#unit-test-development)
  - [Kubernetes](#kubernetes)
    - [Health check](#health-check)
    - [Helm](#helm)
    - [Helmfile](#helmfile)
    - [Skaffold](#skaffold)
- [HTTP Request](#http-request)
- [Trigger Rate Limiting](#trigger-rate-limiting)
- [HTTP Response Status Codes](#http-response-status-codes)
- [Configuration](#configuration)
  - [Minimal](#minimal)
  - [Base path](#base-path)
  - [Blocking](#blocking)
  - [Path replacing (URL rewriting)](#path-replacing-url-rewriting)
  - [Rate Limiting](#rate-limiting)
  - [Token Validation](#token-validation)
    - [JSON Web Key Set (JWKS)](#json-web-key-set-jwks)
    - [Static public key](#static-public-key)
  - [Environment Variables](#environment-variables)
  - [Request Matching](#request-matching)
    - [Path Prefix](#path-prefix)
      - [Catch-all](#catch-all)
    - [HTTP Method](#http-method)
    - [Hostname](#hostname)
      - [Exact match](#exact-match)
      - [Subdomain (wildcard matching)](#subdomain-wildcard-matching)
    - [HTTP Headers](#http-headers)
- [Artifacts](#artifacts)
  - [Binary](#binary)
  - [Docker Image](#docker-image)
  - [Helm Chart](#helm-chart)
- [Credits](#credits)
- [TODO](#todo)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Dependencies (for development)

* Make
* Docker
* Docker Compose
* HTTPie

# Quick start

Start proxy

```sh
make dev
```

Make HTTP request

```sh
http :8080/health
```

# Development

Any changes made to the go or yml files will result in the application being recompiled and relaunched

```sh
make dev
```

## Unit Test Development

Just like `make dev`, the below command will rerun tests if any changes are made to either go or yml files

```sh
make dev-test
```

## Kubernetes

### Health check

/gateway/health

### Helm

...

### Helmfile

...

### Skaffold

Just like `make dev`, the below will automatically recompiled and relaunched the application, but targeting a Kubernetes
cluster.

```sh
skaffold dev
```

# HTTP Request

Perform a http request using HTTPie

```sh
http :8080
```

# Trigger Rate Limiting

```sh
seq 10 | xargs -P 4 -I '{}' http post :8080/health
```

# HTTP Response Status Codes

The following status codes could be returned by the gateway but could also originate from the server which the request
was proxied to

| Code | Meaning             | Reason                 |
|------|---------------------|------------------------|
| 403  | Forbidden           | Invalid token          |
| 421  | Misdirected Request | No matching rule found |
| 429  | Too Many Requests   | Rate limits exceeded   |

# Configuration

## Minimal

The following configuration will proxy all get requests to /health, to the `defaultBackend`.

All other requests will be rate limited before being proxied to the same backend.

```yml
serverPort: 8080
defaultBackend: backend

backends:
  - name: backend
    url: http://backend:8080

rules:
  - pathPrefix: /health
    method: GET
  - pathPrefix: /
    requestPerSecond: 2
    burst: 2
```

## Base path

It's possible to define a base path property. The value of such will be prepended to all path prefixes.

```yml
basePath: /test
rules:
  - pathPrefix: /health
```

The above configuration will match the incoming request against the path prefix `/test/health`.

## Blocking

Access to certain paths can be blocked by using the `block` property of a rule

```yml
rules:
  - pathPrefix: /prefix
    block: true
```

## Path replacing (URL rewriting)

Multiple backends might serve content from the same paths. To accommodate this use case a URL can be rewritten by
replacing part of its path.

```yml
  - pathPrefix: /something
    pathReplace:
      target: thing
      replacement: body
```

The above rule would rewrite the incoming request path to `/somebody`

The `replacement` key is optional and defaults to an empty string.

```yml
  - pathPrefix: /something
    pathReplace:
      target: thing
```

The above rule would rewrite the incoming request path to `/some`

## Rate Limiting

Rate limiting is done per rule basis by defining the number of requests which are allowed per second and how much burst
to allowed.

```yml
rules:
  - pathPrefix: /
    requestPerSecond: 2
    burst: 2
```

This is implemented using [tollbooth](https://github.com/didip/tollbooth).

## Token Validation

Token validation is supported via public key validation.

> **_NOTE:_** The `jwks` and `jwt` configurations are mutually exclusive and if both are defined the `jwks` configuration is evaluated first and takes precedence.

### JSON Web Key Set (JWKS)

A public key can be retrieved from a remote JWKS if a configuration like the below is specified

```yml
authentication:
  jwks:
    host: http://jwks/jwks.json
    index: 0
    minimumRefreshInterval: 960 # 15 minutes
```

### Static public key

A static public key can be configured as shown below.

```yml
authentication:
  jwt:
    publicKey: |
      -----BEGIN PUBLIC KEY-----
      MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAtYrBsSkVGXZKQL13lbmd
      xFCQcvi6KIssjz3KOHIko/Da6sxE2w67OL84t98wCYbmIuq6xTK6qpEqEs1LaqQS
      DnCs2VNDTLk4D1J42R63OpJQfOfebzhTJLx6KldyK2FRGXWILY7AzcoqyuLk433s
      lHk6/yFDYgBA4COofeXZvXtUazuzpBWTZCxpEh341ob6XQ5juLYrqr/80XLYzXiu
      N1iz24ulxSnD0GV4cRfHEnnzN3oYFzoYTcTQB6dffNAs/ADHNA9IemyLbT0ugvbf
      L5MOEBOftYLRwmGFWrXf5s9jccku0FPid2wtZEwsv5Sa+Yvr36KHtrr+PSFksOB1
      0QIDAQAB
      -----END PUBLIC KEY-----
```

The incoming request needs to specify a token in the HTTP authorization header.

It can be either as

> Authorization: token

Or prefixed with the string "Bearer"

> Authorization: Bearer token

An example of an authorizing request can be seen below

```sh
export ACCESS_TOKEN=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQ3OTExNTQ0MTUsImlhdCI6MTYzNzU1NDQxNX0.PtQp6_k5bQ9KE9uk520i4emVnUmxFD8DxyeZsfzgT6CY2oMyXEm7zlIA-4_xz2Q7CrSeqnWxpy0coK9MN0EPE2vhFomTrP6D3l7_lX6Dyn1gH6zWpjC_dRqOSRv3AqS3buZiC-vNwCatLhu6WE74cykBAE2veIr8Gp_ebiITXJKiHBNaTlPk2WEfcJ1NL3g7nafy6l-V4h2-Vj3tapJQiLfpgReIXYIswFYH7En7qy94fL0eOUbZzQI9fOuiXvAN-owR3GYcbwz9Hll23VACWsekMJdDBEgUSdek9JOmRHGxko6FE79-_ClYvF1dGUgZB2mDwY_xF2TOG2q3XDi9Aw
http get :8080/ "Authorization: Bearer $ACCESS_TOKEN"
```

## Environment Variables

The following environment variables are support and will take precedence over values defined in the configuration file.

| Variable             | Configuration path           |
|----------------------|------------------------------|
| APIG_SERVER_PORT     | serverPort                   |
| APIG_BASE_PATH       | basePath                     |
| APIG_DEFAULT_BACKEND | defaultBackend               |
| APIG_PUBLIC_KEY      | authentication.jwt.publicKey |

## Request Matching

A request will first be matched by HTTP method and path prefix. If a match is found HTTP headers will be evaluated.

```yml
  - pathPrefix: /health
    method: GET
    headers:
      User-Agent:
        - HTTPie/2.6.0
```

### Path Prefix

The pathPrefix property is the only mandatory rule property

A rule as simple as the following will match all incoming requests with a path starting with /health

```yml
- pathPrefix: /health
```

#### Catch-all

The below rule serves as a catch-all rule

```yml
- pathPrefix: / # catch all
```

A catch-all rule supports the normal HTTP method and http header matching properties. If any of such are defined it'll *
only* catch requests matching those criteria.

### HTTP Method

A tree entry will be created for each HTTP method and the path prefix configured. The entry will be a concatenation of
the HTTP method and the path prefix.

> method+pathPrefix

If the method property isn't defined, an entry for each of the following methods will be created.

* GET
* HEAD
* POST
* PUT
* PATCH
* DELETE
* CONNECT
* OPTIONS
* TRACE

It's therefore considered best practice to specify the HTTP method on a rule if possible.

### Hostname

#### Exact match

Requests can be matched against an exact hostname as seen below.

```yml
- pathPrefix: /
  hostname: domain.org
```

#### Subdomain (wildcard matching)

If the hostname is prefixed with a "*". Any subdomain of `domain.org` will match the below rule.

```yml
- pathPrefix: /
  hostname: "*.domain.org"
```

> **_NOTE:_** The "*" is considered a wildcard for subdomains only and can only be used as a prefix. The following will not work `doma*.org`.

### HTTP Headers

```yml
- pathPrefix: /
  headers:
    User-Agent:
      - HTTPie/2.6.0
```

# Artifacts

The following artifact can be produced by invoking make targets

## Binary

```sh
make binary
```

## Docker Image

The docker image can be created using the below command

```sh
make docker-image
```

The below command will push the image to the registry defined under the `prod` service in the docker compose file.

```sh
make push-docker-image
```

## Helm Chart

The helm chart can be packed using the below command

```sh
make helm-chart
```

The below command will publish the chart

```sh
make publish-helm
```

> **_NOTE:_** The environment variables `CHART_AUTH_USER` and `CHART_AUTH_PASS` needs to contain valid credentials

# Credits

* https://github.com/lestrrat-go/jwx
* https://github.com/hashicorp/go-immutable-radix
* https://github.com/didip/tollbooth

# TODO

* Comment exported types and functions
* Config validation
* * If publicKey == nil && c.Rules... wants to authenticate... complain
* * If config.DefaultBackend == nil && c.Rules... has no backend... complain
* * Identical rules
* Don't return 502 when the backend isn't found
2022/01/04 01:54:11 http: proxy error: dial tcp: lookup im-job-dev.instance-manager-dev.svc on 172.20.0.10:53: no such host
