// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"context"

	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/manager"
	"openpitrix.io/notification/pkg/pb"
)

func (s *Server) Checker(ctx context.Context, req interface{}) error {
	switch r := req.(type) {
	case *pb.ServiceConfig:
		return manager.NewChecker(ctx, r).
			Required(constants.ServiceCfgProtocol, constants.ServiceCfgEmailHost, constants.ServiceCfgPort, constants.ServiceCfgDisplayEmail, constants.ServiceCfgEmail, constants.ServiceCfgPassword).
			Exec()

	}
	return nil
}
