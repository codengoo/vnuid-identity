# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Install dependencies
COPY . .
RUN go mod download

# Build source files
RUN go build -o server

# Deployment
FROM alpine:latest
WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/server .

# Expose the application port
# ENV DATABASE_URL=postgres://examuser:exam1234@localhost:5432/vnuid?sslmode=disable
# ENV JWT_TOKEN=your-secret-key
# ENV PORT=3333
# EXPOSE 3333

# Run the server
CMD ["./server"]