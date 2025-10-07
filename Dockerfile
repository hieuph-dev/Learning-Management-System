# Stage 1: Build stage
FROM golang:1.24.1-alpine AS builder

# Cài đặt các dependencies cần thiết
RUN apk add --no-cache git gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Stage 2: Final stage
FROM alpine:latest

# Cài đặt ca-certificates cho HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Set timezone
ENV TZ=Asia/Ho_Chi_Minh

WORKDIR /root/

# Copy binary từ builder stage
COPY --from=builder /app/main .

# Copy .env file (optional, nên dùng docker-compose env)
# COPY .env .

# Tạo thư mục uploads
RUN mkdir -p /root/uploads

# Expose port
EXPOSE 8080

# Run application
CMD ["./main"]