FROM resin/raspberrypi3-alpine-golang:slim
RUN apk add --no-cache --update g++ make autoconf automake libtool nasm wget git

WORKDIR /

RUN wget https://github.com/mozilla/mozjpeg/releases/download/v3.2-pre/mozjpeg-3.2-release-source.tar.gz && \
    tar -xvzf mozjpeg-3.2-release-source.tar.gz && \
    rm mozjpeg-3.2-release-source.tar.gz && \
    cd mozjpeg && \
    ./configure && \
    make install && \
    cd / && rm -rf mozjpeg && \
    ln -s /opt/mozjpeg/bin/jpegtran /usr/local/bin/jpegtran && \
    ln -s /opt/mozjpeg/bin/cjpeg /usr/local/bin/cjpeg

RUN go get -v github.com/nickalie/go-binwrapper && \
    go get -v github.com/stretchr/testify/assert

RUN mkdir -p $GOPATH/src/github.com/nickalie/go-mozjpegbin
COPY . $GOPATH/src/github.com/nickalie/go-mozjpegbin
WORKDIR $GOPATH/src/github.com/nickalie/go-mozjpegbin
RUN go test -v ./...