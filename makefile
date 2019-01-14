# Copyright 2018 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

TARG.Name:=notification
TRAG.Gopkg:=openpitrix.io/notification
TRAG.Version:=openpitrix.io/pkg/version


DOCKER_TAGS=latest
BUILDER_IMAGE=openpitrix/openpitrix-builder:release-v0.2.3
RUN_IN_DOCKER:=docker run -it -v `pwd`:/go/src/$(TRAG.Gopkg) -v `pwd`/tmp/cache:/root/.cache/go-build  -w /go/src/$(TRAG.Gopkg) -e GOBIN=/go/src/$(TRAG.Gopkg)/tmp/bin -e USER_ID=`id -u` -e GROUP_ID=`id -g` $(BUILDER_IMAGE)

# Get project build flags
define get_build_flags
    $(eval SHORT_VERSION=$(shell git describe --tags --always --dirty="-dev"))
    $(eval SHA1_VERSION=$(shell git show --quiet --pretty=format:%H))
	$(eval DATE=$(shell date +'%Y-%m-%dT%H:%M:%S'))
	$(eval BUILD_FLAG= -X $(TRAG.Version).ShortVersion="$(SHORT_VERSION)" \
		-X $(TRAG.Version).GitSha1Version="$(SHA1_VERSION)" \
		-X $(TRAG.Version).BuildDate="$(DATE)")
endef


.PHONY: generate-in-local
generate-in-local: ## Generate code from protobuf file in local
	cd ./api && make generate

.PHONY: generate
generate: generate-in-local ## Generate code from protobuf file in docker
	@echo "generate done"



.PHONY: build-flyway
build-flyway: ## Build custom flyway image
	docker build -t $(TARG.Name):flyway -f ./pkg/db/Dockerfile ./pkg/db/


.PHONY: build
build:# fmt build-flyway ## Build all openpitrix images
    #cd ./cmd/server
    #go build
	mkdir -p ./tmp/bin
	$(call get_build_flags)
	$(RUN_IN_DOCKER) time go install -tags netgo -v -ldflags '$(BUILD_FLAG)' $(foreach cmd,$(CMDS),$(TRAG.Gopkg)/cmd/$(cmd))
	docker image prune -f 1>/dev/null 2>&1
	@echo "build done"