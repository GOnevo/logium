VERSION 0.6
FROM golang:1.18-alpine

WORKDIR /app

ENV CGO_ENABLED 0

build:
    COPY . .
    RUN go build -ldflags="-s -w" -o logium main.go
    SAVE ARTIFACT logium /app AS LOCAL bin/logium
