#
# Build image: docker build -t irita/nftp .
#
FROM golang:1.17.5-alpine3.15 as builder

# Set up dependencies
ENV PACKAGES make gcc git

WORKDIR $GOPATH/src
COPY . .

# Install minimum necessary dependencies, build binary
RUN apk add --no-cache $PACKAGES && \
    make install

FROM alpine:3.12

COPY --from=builder /go/bin/nftp /usr/local/bin/nftp
CMD ["sh","-c","nftp start"]
