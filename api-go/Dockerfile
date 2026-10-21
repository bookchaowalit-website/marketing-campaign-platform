# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git make

# Copy source code first (so we can run go mod tidy properly)
COPY . .

# Download all dependencies and tidy
RUN go mod tidy
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata \
    && addgroup -S app && adduser -S -G app app

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Never run as root
USER app

# Expose ports
EXPOSE 8080 9090

# Run the application
CMD ["./main"]
