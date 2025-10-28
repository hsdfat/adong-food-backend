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
FROM alpine:latest

# Set timezone
ENV TZ=Asia/Ho_Chi_Minh

# Create app user for security (non-root)
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder --chown=appuser:appuser /tmp/main .

# Copy .env file if exists (optional - can use env vars instead)
COPY --from=builder --chown=appuser:appuser /app/.env* ./

# Set ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose application port
EXPOSE 18080

# Health check
HEALTHCHECK --interval=30s \
    --timeout=10s \
    --start-period=5s \
    --retries=3 \
    CMD curl -f http://localhost:18080/health || exit 1

# Run the application
CMD ["./main"]