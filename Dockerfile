# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies for SQLite
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application with CGO enabled
RUN CGO_ENABLED=1 go build -o main ./cmd/bot

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies for SQLite
RUN apk add --no-cache sqlite-libs

# Copy the binary from builder
COPY --from=builder /app/main .

# Define entrypoint with flags
ENTRYPOINT ["./main", "-token"]