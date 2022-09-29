# syntax=docker/dockerfile:1
FROM golang:1.18-alpine

WORKDIR /app

RUN apk add build-base

COPY go.mod ./
COPY go.sum ./
COPY ./components/* ./components/
COPY main.go ./

RUN go build -o bot

CMD [ "./bot" ]
