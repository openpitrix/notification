# Copyright 2018 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

FROM golang:1.11-alpine3.7 as builder

# install tools
RUN apk add --no-cache git

# install /usr/bin/nsenter
RUN apk add --no-cache util-linux

WORKDIR /go/src/openpitrix.io//notification
COPY . .

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

#to fix a issue: verifying ...: checksum mismatch
RUN rm go.sum;go mod download
RUN mkdir -p /openpitrix_bin
RUN go build -v -a -installsuffix cgo -ldflags '-w' -o /openpitrix_bin/notification-manager cmd/server/main.go

FROM alpine:3.7
COPY --from=builder /openpitrix_bin/notification-manager /usr/local/bin/
EXPOSE 9201
CMD ["/usr/local/bin/notification-manager"]
