# ============================================
# Stage 1: Build Stage
# ============================================
FROM golang:1.25-alpine AS builder
WORKDIR /app

RUN --mount=type=bind,target=/app \
    go build -o /tmp/main cmd/main.go

# ============================================
# Stage 2: Runtime Stage
# ============================================
FROM scratch

# Copy CA certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime
ENV TZ=Asia/Ho_Chi_Minh

# Copy binary from builder stage
COPY --from=builder /tmp/main /main

# Use non-root user (numeric UID for scratch)
USER 1000:1000

# Expose application port
EXPOSE 18080

# Run the application
CMD ["/main"]