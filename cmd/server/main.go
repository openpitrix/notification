// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package main

import (
	"log"

	"openpitrix.io/notification/pkg/services/notification"
)

func main() {
	log.Println("Starting server...")

	notification.Serve()

	log.Println("Server shuting down...")

}
