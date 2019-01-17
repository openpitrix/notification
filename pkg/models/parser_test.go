// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"strings"
	"testing"
)

func TestGenTasksFromJob(t *testing.T) {
	emailsArray := strings.Split("johuo@yunify.com;danma@yunify.com", ";")
	for _, email := range emailsArray {
		println(email)
	}
}
