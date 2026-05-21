# Stage 1: Build
FROM golang:1.25-alpine AS builder

ARG GOPROXY=https://goproxy.cn,direct
ENV GOPROXY=$GOPROXY

WORKDIR /app

# 依赖缓存层
COPY go.mod go.sum ./
RUN go mod download

# 构建
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Stage 2: Runtime
FROM alpine:3.21

WORKDIR /app

# 时区 & CA 证书
RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/server .
COPY config.docker.yaml ./config.yaml
COPY docker-entrypoint.sh .
RUN chmod +x docker-entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["./docker-entrypoint.sh"]
