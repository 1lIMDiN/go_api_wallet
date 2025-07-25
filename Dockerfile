FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o go_api_wallet
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/go_api_wallet .
COPY --from=builder /app/config.env .

EXPOSE 8080
CMD ["./go_api_wallet"]