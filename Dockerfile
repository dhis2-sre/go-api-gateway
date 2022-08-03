FROM golang:1.18-alpine AS build
ARG REFLEX_VERSION=v0.3.1
RUN apk add gcc musl-dev git
WORKDIR /src
RUN go install github.com/cespare/reflex@${REFLEX_VERSION}
COPY go.mod go.sum ./
RUN go mod download -x
COPY . .
#RUN go build -o /app/go-api-gateway -ldflags "-s -w" ./cmd/serve
RUN go build -o /app/go-api-gateway -ldflags "-s -w" ./cmd/mux

FROM alpine:3.15
RUN apk --no-cache -U upgrade
WORKDIR /app
COPY --from=build /app/go-api-gateway .
USER guest
ENTRYPOINT ["/app/go-api-gateway"]
