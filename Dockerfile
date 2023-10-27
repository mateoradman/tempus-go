FROM golang:1.21-alpine3.16 AS builder

WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY config.env .
COPY start.sh .
COPY wait-for.sh .
COPY internal/db/migration ./internal/db/migration
