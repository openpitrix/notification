// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package emailutil

import (
	"errors"
	"net/smtp"
)

type noStartTLSPlainAuth struct {
	identity string
	username string
	password string
	host     string
}

func (a *noStartTLSPlainAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	//compared net/smtp  Auth.go  PlainAuth,remove below lines.
	//if !server.TLS && !isLocalhost(server.Name) {
	//	return "", nil, errors.New("unencrypted connection")
	//}

	if server.Name != a.host {
		return "", nil, errors.New("wrong host name")
	}
	resp := []byte(a.identity + "\x00" + a.username + "\x00" + a.password)
	return "PLAIN", resp, nil
}

func (a *noStartTLSPlainAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		// We've already sent everything.
		return nil, errors.New("unexpected server challenge")
	}
	return nil, nil
}
