// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"context"
	"fmt"

	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/gerr"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/pb"
	rs "openpitrix.io/notification/pkg/services/notification/resource_control"
	"openpitrix.io/notification/pkg/util/pbutil"
)

func (s *Server) SetServiceConfig(ctx context.Context, req *pb.ServiceConfig) (*pb.SetServiceConfigResponse, error) {
	err := ValidateSetServiceConfigParams(ctx, req)
	if err != nil {
		return nil, err
	}

	rs.SetServiceConfig(req)
	return &pb.SetServiceConfigResponse{
		IsSucc: pbutil.ToProtoBool(true),
	}, nil

}

func (s *Server) GetServiceConfig(ctx context.Context, req *pb.GetServiceConfigRequest) (*pb.ServiceConfig, error) {
	var ServiceTypes = []string{
		constants.ServiceTypeEmail,
	}

	serviceTypes := req.GetServiceType()
	if len(serviceTypes) == 0 {
		serviceTypes = ServiceTypes
	}

	var emailCfg *pb.EmailServiceConfig

	scCfg := &pb.ServiceConfig{}

	for _, scType := range serviceTypes {
		if scType == constants.NotifyTypeEmail {
			emailCfg = rs.GetEmailServiceConfig()
			break
		}
	}
	if emailCfg == nil {
		err := gerr.NewWithDetail(ctx, gerr.Internal, fmt.Errorf("Can not get EmailServiceConfig"), gerr.ErrorGetServiceConfigFailed)
		return nil, err
	}
	logger.Infof(ctx, "Get ServiceConfig successfully. %+v", emailCfg)

	scCfg.EmailServiceConfig = emailCfg
	return scCfg, nil
}

func (s *Server) CreateNotification(ctx context.Context, req *pb.CreateNotificationRequest) (*pb.CreateNotificationResponse, error) {
	notification := models.NewNotification(
		req.GetContentType().GetValue(),
		req.GetTitle().GetValue(),
		req.GetContent().GetValue(),
		req.GetShortContent().GetValue(),
		req.GetAddressInfo().GetValue(),
		req.GetOwner().GetValue(),
		req.GetExpiredDays().GetValue(),
	)

	err := rs.RegisterNotification(ctx, notification)
	if err != nil {
		return nil, err
	}

	tasks, err := SplitNotificationIntoTasks(notification)
	if err != nil {
		logger.Errorf(ctx, "Split notification into tasks failed, [%+v]", err)
		return nil, err
	}

	for _, task := range tasks {
		err = rs.RegisterTask(ctx, task)
		if err != nil {
			return nil, err
		}
		err = s.controller.taskQueue.Enqueue(task.TaskId)
		if err != nil {
			logger.Errorf(ctx, "Push task [%s] into etcd failed, [%+v]", task.TaskId, err)
			return nil, err
		}
	}

	// Enqueue notification after tasks.
	err = s.controller.notificationQueue.Enqueue(notification.NotificationId)
	if err != nil {
		logger.Errorf(ctx, "Push notification [%s] into etcd failed, [%+v]", notification.NotificationId, err)
		return nil, err
	}

	return &pb.CreateNotificationResponse{
		NotificationId: pbutil.ToProtoString(notification.NotificationId),
	}, nil
}

func (s *Server) DescribeNotifications(ctx context.Context, req *pb.DescribeNotificationsRequest) (*pb.DescribeNotificationsResponse, error) {
	nfs, nfCnt, err := rs.DescribeNotifications(ctx, req)
	if err != nil {
		logger.Errorf(ctx, "Failed to describe notifications, error: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	nfPbSet := models.ParseNfSet2PbSet(nfs)
	res := &pb.DescribeNotificationsResponse{
		TotalCount:      uint32(nfCnt),
		NotificationSet: nfPbSet,
	}
	return res, nil
}

func (s *Server) RetryNotifications(ctx context.Context, req *pb.RetryNotificationsRequest) (*pb.RetryNotificationsResponse, error) {
	return &pb.RetryNotificationsResponse{}, nil
}

func (s *Server) DescribeTasks(ctx context.Context, req *pb.DescribeTasksRequest) (*pb.DescribeTasksResponse, error) {
	tasks, taskCnt, err := rs.DescribeTasks(ctx, req)
	if err != nil {
		logger.Errorf(ctx, "Failed to describe tasks, error: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	taskPbSet := models.ParseTaskSet2PbSet(tasks)
	res := &pb.DescribeTasksResponse{
		TotalCount: uint32(taskCnt),
		TaskSet:    taskPbSet,
	}
	return res, nil
}

func (s *Server) RetryTasks(ctx context.Context, req *pb.RetryTasksRequest) (*pb.RetryTasksResponse, error) {
	return &pb.RetryTasksResponse{}, nil
}

func (s *Server) CreateAddress(ctx context.Context, req *pb.CreateAddressRequest) (*pb.CreateAddressResponse, error) {
	addr := models.NewAddress(req)
	addrId, err := rs.CreateAddress(ctx, addr)

	if err != nil {
		logger.Errorf(ctx, "Failed to create address [%+v], %+v", addr, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}
	logger.Debugf(ctx, "Create new address, [%+v]", addr)

	return &pb.CreateAddressResponse{
		AddressId: pbutil.ToProtoString(addrId),
	}, nil
}

func (s *Server) DescribeAddresses(ctx context.Context, req *pb.DescribeAddressesRequest) (*pb.DescribeAddressesResponse, error) {
	addrs, addrCnt, err := rs.DescribeAddresses(ctx, req)
	if err != nil {
		logger.Errorf(ctx, "Failed to Describe Addresses, error: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	addressPbSet := models.ParseAddressSet2PbSet(addrs)
	res := &pb.DescribeAddressesResponse{
		TotalCount: uint32(addrCnt),
		AddressSet: addressPbSet,
	}
	return res, nil
}

func (s *Server) ModifyAddress(ctx context.Context, req *pb.ModifyAddressRequest) (*pb.ModifyAddressResponse, error) {
	addressId, err := rs.ModifyAddress(ctx, req)
	if err != nil {
		logger.Errorf(ctx, "Failed to Modify Address [%s], %+v", addressId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed, addressId)
	}
	return &pb.ModifyAddressResponse{
		AddressId: pbutil.ToProtoString(addressId),
	}, nil

}

func (s *Server) DeleteAddresses(ctx context.Context, req *pb.DeleteAddressesRequest) (*pb.DeleteAddressesResponse, error) {
	return &pb.DeleteAddressesResponse{}, nil
}

func (s *Server) CreateAddressList(ctx context.Context, req *pb.CreateAddressListRequest) (*pb.CreateAddressListResponse, error) {
	return &pb.CreateAddressListResponse{}, nil
}

func (s *Server) DescribeAddressList(ctx context.Context, req *pb.DescribeAddressListRequest) (*pb.DescribeAddressListResponse, error) {
	return &pb.DescribeAddressListResponse{}, nil
}

func (s *Server) ModifyAddressList(ctx context.Context, req *pb.ModifyAddressListRequest) (*pb.ModifyAddressListResponse, error) {
	return &pb.ModifyAddressListResponse{}, nil
}

func (s *Server) DeleteAddressList(ctx context.Context, req *pb.DeleteAddressListRequest) (*pb.DeleteAddressListResponse, error) {
	return &pb.DeleteAddressListResponse{}, nil
}
