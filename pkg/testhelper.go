// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import "flag"

var (
	LocalDevEnvEnabled = flag.Bool("LocalDevEnvEnabled", false, "disenable Local Dev Env setting")
	//LocalDevEnvEnabled = flag.Bool("LocalDevEnvEnabled", true, "enable Local Dev Env setting")
)
