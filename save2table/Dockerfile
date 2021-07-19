FROM golang:alpine AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 移动到工作目录：/build
WORKDIR /build

RUN go env -w GOPROXY="https://goproxy.cn,direct"

# 复制项目中的 go.mod 和 go.sum文件并下载依赖信息
COPY go.mod .
COPY go.sum .
RUN go mod download

# 将代码复制到容器中
COPY . .

# 将代码编译成二进制可执行文件app
RUN go build -o app .

###################
# 接下来创建一个小镜像
###################
FROM debian:stretch-slim

COPY ./conf /conf

COPY --from=builder /build/app /
# 声明服务端口
EXPOSE 8081

# 启动容器时运行的命令
ENTRYPOINT ["./app"]