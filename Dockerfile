FROM golang:1.14.7-alpine

ENV GOPROXY https://goproxy.cn
ENV GO111MODULE on
ENV CGO_ENABLED 0

WORKDIR /go/src/ChatRoom

COPY . .

RUN go mod download