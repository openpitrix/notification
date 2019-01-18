// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"testing"
	"time"
)

func TestServe(t *testing.T) {
	c := NewController()

	go c.ExtractTasks()
	go c.HandleTask("A")
	go c.HandleTask("B")
	//
	for {
		//println("...")
		time.Sleep(2 * time.Second)
	}
}
