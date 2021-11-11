FROM golang:1.16-alpine AS build
RUN apk add gcc musl-dev
WORKDIR /src
RUN go get github.com/cespare/reflex
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/go-rate-limit -ldflags "-s -w" ./cmd/go-rate-limit

FROM alpine:3.14
RUN apk --no-cache -U upgrade
WORKDIR /app
COPY --from=build /app/go-rate-limit .
USER guest
CMD ["/app/go-rate-limit"]
