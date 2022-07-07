FROM golang:1.18.1-alpine3.15 AS builder

WORKDIR /build

RUN apk update
RUN apk add --no-cache git

ARG GITHUB_TOKEN

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOARCH="amd64" \
    GOOS=linux

COPY go.mod go.sum ./

RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"

RUN go mod download

COPY . .

# building
RUN go build -o kouda ./cmd/kouda/main.go

FROM alpine:20210804

WORKDIR /app

ENV APP_HOME=/app

RUN mkdir -p ${APP_HOME}/config

COPY --from=builder /build/config/barong.yml ./config/
COPY --from=builder /build/kouda ./
