// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"context"

	"openpitrix.io/notification/pkg/models"

	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/manager"
	"openpitrix.io/notification/pkg/pb"
)

var NfStatuses = []string{
	constants.StatusPending,
	constants.StatusSending,
	constants.StatusSuccessful,
	constants.StatusFailed,
}

var NfTypes = []string{
	constants.NotifyTypeEmail,
	constants.NotifyTypeWeb,
	constants.NotifyTypeMobile,
	constants.NotifyTypeSms,
	constants.NotifyTypeWeChat,
}

var Statuses = []string{
	constants.StatusActive,
	constants.StatusDisabled,
	constants.StatusDeleted,
}

var ContentTypes = []string{
	constants.ContentTypeInvite,
	constants.ContentTypeverify,
	constants.ContentTypeFee,
	constants.ContentTypeBusiness,
	constants.ContentTypeOther,
}

func (s *Server) Checker(ctx context.Context, req interface{}) error {
	switch r := req.(type) {
	case *pb.ServiceConfig:
		return manager.NewChecker(ctx, r).
			Required(constants.ServiceCfgProtocol, constants.ServiceCfgEmailHost, constants.ServiceCfgPort, constants.ServiceCfgDisplayEmail, constants.ServiceCfgEmail, constants.ServiceCfgPassword).
			Exec()
	case *pb.GetServiceConfigRequest:
		return manager.NewChecker(ctx, r).
			StringChosen("service_type", NfTypes).
			Exec()
	case *pb.CreateNotificationRequest:
		return manager.NewChecker(ctx, r).
			Required(models.NfColContentType, models.NfColTitle, models.NfColShortContent, models.NfColAddressInfo).
			StringChosen("content_type", ContentTypes).
			Exec()
	case *pb.DescribeNotificationsRequest:
		return manager.NewChecker(ctx, r).
			StringChosen("status", NfStatuses).
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
			StringChosen("service_type", NfTypes).
			Exec()
	case *pb.DescribeAddressesRequest:
		return manager.NewChecker(ctx, r).
			StringChosen("status", Statuses).
			StringChosen("notify_type", NfTypes).
			Exec()
	case *pb.ModifyAddressRequest:
		return manager.NewChecker(ctx, r).
			Required(models.AddrColId).
			StringChosen("service_type", NfTypes).
			Exec()
	case *pb.DeleteAddressesRequest:
		return manager.NewChecker(ctx, r).
			Required(models.TaskColTaskId).
			Exec()
		//case *pb.CreateAddressListRequest:
		//	return manager.NewChecker(ctx, r).
		//		Required(models.TaskColTaskId).
		//		Exec()

	}

	return nil
}
