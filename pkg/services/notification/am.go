// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"context"

	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/manager"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/pb"
)

func (s *Server) Checker(ctx context.Context, req interface{}) error {
	switch r := req.(type) {
	case *pb.ServiceConfig:
		return manager.NewChecker(ctx, r).
			Required(models.ServiceCfgProtocol, models.ServiceCfgEmailHost, models.ServiceCfgPort, models.ServiceCfgDisplayEmail, models.ServiceCfgEmail, models.ServiceCfgPassword).
			StringChosen(models.ServiceCfgProtocol, models.ProtocolTypes).
			Exec()
	case *pb.GetServiceConfigRequest:
		return manager.NewChecker(ctx, r).
			StringChosen(models.ServiceType, constants.NotifyTypes).
			Exec()
	case *pb.CreateNotificationRequest:
		return manager.NewChecker(ctx, r).
			Required(models.NfColContentType, models.NfColTitle, models.NfColAddressInfo).
			StringChosen(models.NfColContentType, models.ContentTypes).
			Exec()
	case *pb.DescribeNotificationsRequest:
		return manager.NewChecker(ctx, r).
			StringChosen(models.NfColStatus, constants.SendingStatuses).
			Exec()
	case *pb.RetryNotificationsRequest:
		return manager.NewChecker(ctx, r).
			Required(models.NfColId).
			Exec()
	case *pb.RetryTasksRequest:
		return manager.NewChecker(ctx, r).
			Required(models.TaskColTaskId).
			Exec()
	case *pb.CreateAddressRequest:
		return manager.NewChecker(ctx, r).
			Required(models.AddrColAddress, models.AddrColNotifyType).
			StringChosen(models.ServiceType, constants.NotifyTypes).
			Exec()
	case *pb.DescribeAddressesRequest:
		return manager.NewChecker(ctx, r).
			StringChosen(models.AddrColStatus, constants.RecordStatuses).
			StringChosen(models.AddrColNotifyType, constants.NotifyTypes).
			Exec()
	case *pb.ModifyAddressRequest:
		return manager.NewChecker(ctx, r).
			Required(models.AddrColId).
			StringChosen(models.ServiceType, constants.NotifyTypes).
			Exec()
	case *pb.DeleteAddressesRequest:
		return manager.NewChecker(ctx, r).
			Required(models.AddrColId).
			Exec()
	}

	return nil
}
