# Copyright 2019 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

default:
	docker build -t openpitrix/notification-builder:v1.0.0 .
	@echo "ok"

pull:
	docker pull openpitrix/notification-builder:v1.0.0
	@echo "ok"

run:
	docker run --rm -it -v `pwd`:/root openpitrix/notification-builder:v1.0.0

clean:
	@echo "ok"
