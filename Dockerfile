# Stage 1: Build the Go application
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o woc ./cmd

# Stage 2: Create a minimal image to run the application
FROM alpine:latest

RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /app/woc .

USER appuser

EXPOSE 8080

CMD ["./woc"]