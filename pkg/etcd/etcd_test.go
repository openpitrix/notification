// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package etcd

import (
	"log"
	"testing"

	pkg "openpitrix.io/notification/pkg"
)

func TestConnect(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("LocalDevEnv disabled")
	}
	endpoints := []string{"192.168.0.6:12379"}
	//endpoints := []string{"139.198.121.89:12379"}
	//endpoints := []string{"139.198.121.89:52379"}

	prefix := "test"
	e, err := Connect(endpoints, prefix)
	log.Println(e)
	if err != nil {
		t.Fatal(err)
	}
}
