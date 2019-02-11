// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build etcd

package etcdutil

import (
	"fmt"
	"log"
	"os"
	"testing"

	"openpitrix.io/openpitrix/pkg/etcd"
)

func TestConnect(t *testing.T) {
	endpoints := os.Getenv("NOTIFICATION_ETCD_ENDPOINTS")
	//endpoints := []string{"192.168.0.7:2379"}
	prefix := "test"
	e, err := etcd.Connect(endpoints, prefix)
	log.Println(e)
	if err != nil {
		t.Fatal(err)
	}

}

func TestEnqueue(t *testing.T) {
	endpoints := os.Getenv("NOTIFICATION_ETCD_ENDPOINTS")
	//endpoints := []string{"192.168.0.7:2379"}
	prefix := "test"
	e, err := etcd.Connect(endpoints, prefix)
	if err != nil {
		t.Fatal(err)
	}
	queue := e.NewQueue("notification")
	go func() {
		for i := 0; i < 100; i++ {
			err := queue.Enqueue(fmt.Sprintf("%d", i))
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("Push message to queue, worker number [%d]", i)
		}

	}()
	for i := 0; i < 100; i++ {
		n, err := queue.Dequeue()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Got message [%s] from queue, worker number [%d]", n, i)
	}
}
