# Dockerfile for the shred-service

# Stage 1: Build the executable
FROM golang:1.24 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o /bin/shred-service ./cmd/shred-service

# Stage 2: Final runtime image
FROM ubuntu:latest AS runtime
WORKDIR /app
COPY --from=builder /bin/shred-service /app/shred-service

EXPOSE 8088
CMD ["/app/shred-service"]
