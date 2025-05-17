# syntax=docker/dockerfile:1.4

# Builder stage
FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY --link go.mod go.sum ./
RUN go mod download

COPY --link . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/s3-syncd ./cmd/main.go


# Final runtime image
FROM alpine:3.21

WORKDIR /app

# (Optional) Install CA certificates if your binary makes HTTPS requests
RUN apk add --no-cache ca-certificates

COPY --from=builder /usr/local/bin/s3-syncd /usr/local/bin/s3-syncd

CMD ["/usr/local/bin/s3-syncd"]
