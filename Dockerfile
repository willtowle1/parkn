FROM golang:alpine as builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
COPY cmd cmd
COPY internal internal
COPY .env .env