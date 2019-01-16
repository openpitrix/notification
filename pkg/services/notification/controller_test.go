// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"testing"
	"time"

	"openpitrix.io/notification/pkg/services/notification/service/notification"
	"openpitrix.io/notification/pkg/services/notification/service/task"
)

func TestServe(t *testing.T) {
	nfservice := notification.NewService()
	taskservice := task.NewService()

	c := NewController(nfservice, taskservice)

	go c.ExtractTasks()
	go c.HandleTask("A")
	go c.HandleTask("B")
	//
	for {
		//println("...")
		time.Sleep(2 * time.Second)
	}
}
