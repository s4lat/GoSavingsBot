# syntax=docker/dockerfile:1
FROM golang:1.18-alpine

WORKDIR /app

RUN apk add build-base

COPY . .

RUN go mod download && go mod verify
RUN go build main.go

CMD [ "./main" ]
