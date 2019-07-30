# Copyright 2019 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.


FROM golang:1.12.7-alpine3.10 as builder

RUN apk add --no-cache git curl openssl


RUN export GO111MODULE=on; go get github.com/golang/protobuf/protoc-gen-go@v1.3.2
RUN export GO111MODULE=on; go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.9.5
RUN export GO111MODULE=on; go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@v1.9.5
RUN export GO111MODULE=on; go get golang.org/x/tools/cmd/goimports

# swagger-0.13.0
# RUN go get github.com/go-swagger/go-swagger/cmd/swagger
RUN mkdir -p /swagger && cd /swagger \
	&& wget https://github.com/go-swagger/go-swagger/releases/download/0.13.0/swagger_linux_amd64 \
	&& chmod +x swagger_linux_amd64 && mv swagger_linux_amd64 /go/bin/swagger

# the protoc can't run on alpine,
# we only need the protobuf's stdarnd library in the `/protoc/include`.
RUN mkdir -p /protoc && cd /protoc \
	&& wget https://github.com/google/protobuf/releases/download/v3.5.0/protoc-3.5.0-linux-x86_64.zip \
	&& unzip protoc-3.5.0-linux-x86_64.zip

FROM golang:1.12.7-alpine3.10

RUN apk add --no-cache git protobuf make curl openssl jq rsync upx

COPY --from=builder /protoc/include /usr/local/include
COPY --from=builder /go/bin /go/bin


#docker build -t openpitrix/notification-builder:v1.0.0 .
