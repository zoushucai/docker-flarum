FROM golang:alpine AS builder
# FROM golang:1.23.0 AS builder (如果采用这个系统,下面的某些指令会报错)
## 无论采用了那个系统, 最后采用了多阶段构建, 所以最后的体积都是一样的
LABEL stage=gobuilder

ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build
COPY . .
RUN go mod download
RUN go build -ldflags="-s -w" -o /app/main main.go


FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/main /app/main
EXPOSE 9113
CMD ["./main"]
