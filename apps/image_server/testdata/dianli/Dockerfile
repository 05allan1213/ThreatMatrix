FROM golang:alpine AS builder

# 构建可执行文件
ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

WORKDIR /build
ADD . .
RUN go mod tidy
RUN go build -o main

FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/main /app

# 安装 ca-certificates 以支持 HTTPS 请求
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update && \
    apk --no-cache add ca-certificates

CMD ["./main"]