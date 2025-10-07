# =================================
# Stage 1: Build stage
# =================================
FROM golang:1.23-alpine AS builder

# Cài đặt các dependencies cần thiết
RUN apk add --no-cache git gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy toàn bộ source code
COPY . .

# Build application với CGO enabled (cần cho sqlite nếu dùng)
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api/main.go

# =================================
# Stage 2: Final stage (Production)
# =================================
FROM alpine:latest

# Cài đặt ca-certificates cho HTTPS và timezone
RUN apk --no-cache add ca-certificates tzdata

# Set timezone sang Việt Nam
ENV TZ=Asia/Ho_Chi_Minh

# Tạo user non-root để chạy app (best practice)
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /home/appuser

# Copy binary từ builder stage
COPY --from=builder /app/main .

# Tạo thư mục uploads và set permissions
RUN mkdir -p /home/appuser/uploads && \
    chown -R appuser:appgroup /home/appuser

# Switch sang user non-root
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

# Run application
CMD ["./main"]