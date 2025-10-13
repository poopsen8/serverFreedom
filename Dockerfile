FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o user_server ./cmd/server

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/user_server .

COPY config4.json ./
COPY config/config.yaml ./config/config.yaml

EXPOSE 8080

CMD ["./user_server"]
