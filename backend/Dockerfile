FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Install swag for documentation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate swagger docs
RUN swag init

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy generated docs
COPY --from=builder /app/docs ./docs

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./main"]