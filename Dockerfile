
FROM golang:1.24-alpine AS builder


RUN apk add --no-cache git


WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN go build -o main ./cmd


FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .


COPY .env .env


EXPOSE 8080


CMD sleep 5 && ./main
