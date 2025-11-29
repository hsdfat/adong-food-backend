# ============================================
# Stage 1: Build Stage
# ============================================
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Copy entire project for build (needed for embedded SQL files)
COPY . .

# Build the application
RUN go build -o /tmp/main cmd/main.go

# ============================================
# Stage 2: Runtime Stage
# ============================================
FROM scratch

# Copy CA certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy all timezone data (entire zoneinfo directory)
ENV TZ=Asia/Ho_Chi_Minh

# Copy binary from builder stage
COPY --from=builder /tmp/main /main

# Use non-root user (numeric UID for scratch)
USER 1000:1000

# Expose application port
EXPOSE 18080

# Run the application
CMD ["/main"]