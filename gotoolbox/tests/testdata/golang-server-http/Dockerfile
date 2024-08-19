# Build stage
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gotoolbox-service .

# Final stage
FROM alpine:3.14

# Set working directory
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/gotoolbox-service .

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./gotoolbox-service"]
