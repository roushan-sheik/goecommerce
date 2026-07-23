# Build Stage
FROM golang:alpine AS builder

WORKDIR /app

# Install dependencies first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .
# Note: Main package is in cmd/
RUN go build -o goecommerce-app ./cmd/main.go

# Final Run Stage
FROM alpine:latest

WORKDIR /app

# Copy only the compiled binary from the builder stage
COPY --from=builder /app/goecommerce-app .

# Expose port and configure environment variable
ENV PORT=8080
EXPOSE 8080

# Run the binary
CMD ["./goecommerce-app"]