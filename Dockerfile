FROM golang:1.16-alpine AS build
RUN apk add gcc musl-dev
WORKDIR /src
RUN go get github.com/cespare/reflex
COPY go.mod go.sum ./
RUN go mod download -x
COPY . .
RUN go build -o /app/go-api-gateway -ldflags "-s -w" ./cmd/serve

FROM alpine:3.14
RUN apk --no-cache -U upgrade
WORKDIR /app
COPY --from=build /app/go-api-gateway .
USER guest
ENTRYPOINT ["/app/go-api-gateway"]
