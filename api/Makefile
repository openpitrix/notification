# Copyright 2018 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.


PROTO_FILES=$(sort $(wildcard ./*.proto))
PROTOC_INC_PATH=/usr/local/include
GOOGLEAPIS_PATH=third_party/googleapis
GRPC_GATEWAY_PATH=third_party/github.com/grpc-ecosystem/grpc-gateway/

PROTOC_FLAGS:=-I. -I$(GOOGLEAPIS_PATH) -I$(GRPC_GATEWAY_PATH) -I$(PROTOC_INC_PATH)

GOPATH:=$(shell go env GOPATH)
PWD:=$(shell pwd)

generate: $(PROTO_FILES) Makefile
	@rm -rf ../pkg/pb/*

    #To generate notification.pb.go file
	protoc $(PROTOC_FLAGS) --go_out=plugins=grpc:../pkg/pb ${PROTO_FILES}

	#To generate notification.pb.gw.go file
	protoc $(PROTOC_FLAGS) --grpc-gateway_out=logtostderr=true,allow_delete_body=true:../pkg/pb ${PROTO_FILES}

	#generate notification.pb.go file
	protoc $(PROTOC_FLAGS) --swagger_out=logtostderr=true,allow_delete_body=true:../pkg/apigateway/spec ${PROTO_FILES}

	cd ../pkg/apigateway/spec/ && cat ./*.swagger.json 0.swagger.json | jq --slurp 'reduce .[] as $$item ({}; . * $$item)' | jq -f ./preprocess.jq > ./api.swagger.json
	cd ../pkg/apigateway/spec/ && find . | grep .swagger.json | grep -v "api" | xargs rm

	make -C ../pkg/apigateway/spec
	make -C ../pkg/apigateway/swagger-ui

	goimports -l -w -e -local=openpitrix ../pkg/pb


clean:
	rm -rf $(GOPATH)/src/openpitrix.io/iam/pkg/pb


