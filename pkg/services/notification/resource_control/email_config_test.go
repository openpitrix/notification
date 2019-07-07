// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.
package resource_control

import (
	"testing"

	"openpitrix.io/logger"

	pkg "openpitrix.io/notification/pkg"
	"openpitrix.io/notification/pkg/config"
)

func TestResetEmailCfg(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}
	cfg := config.GetInstance()
	err := ResetEmailCfg(cfg)

	if err != nil {
		logger.Errorf(nil, "Failed to reset email config with data in db, %+v.", err)
	}
	logger.Debugf(nil, "Reset email config with data in db successfully.")

}
