// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/gerr"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/pbutil"
)

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

	err := RegisterNotification(ctx, notification)
	if err != nil {
		return nil, err
	}

	tasks, err := SplitNotificationIntoTasks(notification)
	if err != nil {
		logger.Errorf(ctx, "Split notification into tasks failed, [%+v]", err)
		return nil, err
	}

	for _, task := range tasks {
		err = RegisterTask(ctx, task)
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
	return &pb.DescribeNotificationsResponse{}, nil
}

func (s *Server) RetryNotifications(ctx context.Context, req *pb.RetryNotificationsRequest) (*pb.RetryNotificationsResponse, error) {
	return &pb.RetryNotificationsResponse{}, nil
}

func (s *Server) DescribeTasks(ctx context.Context, req *pb.DescribeTasksRequest) (*pb.DescribeTasksResponse, error) {
	return &pb.DescribeTasksResponse{}, nil
}

func (s *Server) RetryTasks(ctx context.Context, req *pb.RetryTasksRequest) (*pb.RetryTasksResponse, error) {
	return &pb.RetryTasksResponse{}, nil
}

func (s *Server) CreateAddress(ctx context.Context, req *pb.CreateAddressRequest) (*pb.CreateAddressResponse, error) {
	return &pb.CreateAddressResponse{}, nil
}

func (s *Server) DescribeAddresses(ctx context.Context, req *pb.DescribeAddressesRequest) (*pb.DescribeAddressesResponse, error) {
	return &pb.DescribeAddressesResponse{}, nil
}

func (s *Server) ModifyAddress(ctx context.Context, req *pb.ModifyAddressRequest) (*pb.ModifyAddressResponse, error) {
	return &pb.ModifyAddressResponse{}, nil
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

func (s *Server) SetServiceConfig(ctx context.Context, req *pb.ServiceConfig) (*pb.SetServiceConfigResponse, error) {
	err := s.validateSetServiceConfigParams(ctx, req)
	if err != nil {
		return nil, err
	}

	protocol := req.GetEmailServiceConfig().GetProtocol().GetValue()
	emailHost := req.GetEmailServiceConfig().GetEmailHost().GetValue()
	port := req.GetEmailServiceConfig().GetPort().GetValue()
	displayEmail := req.GetEmailServiceConfig().GetDisplayEmail().GetValue()
	email := req.GetEmailServiceConfig().GetEmail().GetValue()
	password := req.GetEmailServiceConfig().GetPassword().GetValue()
	sslEnable := req.GetEmailServiceConfig().GetSslEnable().GetValue()

	os.Setenv("NOTIFICATION_EMAIL_PROTOCOL", protocol)
	os.Setenv("NOTIFICATION_EMAIL_EMAIL_HOST", emailHost)
	os.Setenv("NOTIFICATION_EMAIL_PORT", port)
	os.Setenv("NOTIFICATION_EMAIL_DISPLAY_EMAIL", displayEmail)
	os.Setenv("NOTIFICATION_EMAIL_EMAIL", email)
	os.Setenv("NOTIFICATION_EMAIL_PASSWORD", password)
	os.Setenv("NOTIFICATION_EMAIL_SSL_ENABLE", strconv.FormatBool(sslEnable))

	config.GetInstance().LoadConf()

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
			emailCfg = GetEmailServiceConfig(ctx)
			break
		}
	}
	if emailCfg == nil {
		err := gerr.NewWithDetail(ctx, gerr.Internal, fmt.Errorf("Can not get EmailServiceConfig"), gerr.ErrorGetServiceConfigFailed)
		return nil, err
	}

	scCfg.EmailServiceConfig = emailCfg
	return scCfg, nil
}

func (s *Server) validateSetServiceConfigParams(ctx context.Context, req *pb.ServiceConfig) error {
	email := req.GetEmailServiceConfig().GetEmail().GetValue()
	err := VerifyEmailFmt(ctx, email)
	if err != nil {
		logger.Errorf(ctx, "Failed to validateSetServiceConfigParams [%s], %+v", email, err)
		return err
	}

	displayEmail := req.GetEmailServiceConfig().GetDisplayEmail().GetValue()
	err = VerifyEmailFmt(ctx, displayEmail)
	if err != nil {
		logger.Errorf(ctx, "Failed to validateSetServiceConfigParams [%s], %+v", displayEmail, err)
		return err
	}

	portStr := req.GetEmailServiceConfig().GetPort().GetValue()
	portNum, err := strconv.ParseInt(portStr, 10, 64)
	err = VerifyPortFmt(ctx, portNum)
	if err != nil {
		logger.Errorf(ctx, "Failed to validateSetServiceConfigParams [%s], %+v", portStr, err)
		return err
	}
	return nil
}
