// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package global

import (
	"testing"

	pkg "openpitrix.io/notification/pkg"
)

func TestSetGlobalCfg(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}
	GetInstance()
}
