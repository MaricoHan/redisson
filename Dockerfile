#
# Build image: docker build -t irita/nftp .
#
FROM golang:1.17.5-alpine3.15 as builder

# Set up dependencies
ENV PACKAGES make git libc-dev bash gcc
ENV GOPROXY http://192.168.0.60:8081/repository/go-bianjie/,http://nexus.bianjie.ai/repository/golang-group,https://goproxy.cn,direct
ENV GOPRIVATE gitlab.bianjie.ai
ENV APKPROXY http://192.168.0.60:8081/repository/apk-ustc

WORKDIR $GOPATH/src
COPY . .

# Install minimum necessary dependencies, build binary
RUN sed -i "s+https://dl-cdn.alpinelinux.org/alpine+${APKPROXY}+g" /etc/apk/repositories && \
    apk add --no-cache $PACKAGES && \
    git config --global url."https://GITUSER:GITPASS@gitlab.bianjie.ai".insteadOf "https://gitlab.bianjie.ai" && \
    go mod tidy && \
    make install

FROM alpine:3.15

COPY --from=builder /go/bin/open-api /usr/local/bin/open-api
CMD ["sh","-c","open-api start"]
