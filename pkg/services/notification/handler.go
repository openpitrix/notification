// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"context"
	"crypto/tls"
	"fmt"

	gomail "gopkg.in/gomail.v2"

	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/gerr"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/pb"
	rs "openpitrix.io/notification/pkg/services/notification/resource_control"
	"openpitrix.io/notification/pkg/util/pbutil"
	"openpitrix.io/notification/pkg/util/stringutil"
)

func (s *Server) SetServiceConfig(ctx context.Context, req *pb.ServiceConfig) (*pb.SetServiceConfigResponse, error) {
	err := ValidateSetServiceConfigParams(ctx, req)
	if err != nil {
		return nil, err
	}

	rs.SetServiceConfig(req)
	logger.Debugf(ctx, "Set service config successfully, %+v.", req)

	return &pb.SetServiceConfigResponse{
		IsSucc: pbutil.ToProtoBool(true),
	}, nil

}

func (s *Server) GetServiceConfig(ctx context.Context, req *pb.GetServiceConfigRequest) (*pb.ServiceConfig, error) {
	var ServiceTypes = []string{
		constants.NotifyTypeEmail,
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
		logger.Errorf(ctx, "Failed to get service config, emailserviceconfig.")
		err := gerr.NewWithDetail(ctx, gerr.Internal, fmt.Errorf("Failed to get service config, emailserviceconfig."), gerr.ErrorGetServiceConfigFailed)
		return nil, err
	}
	logger.Debugf(ctx, "Get service config [%+v] successfully.", emailCfg)

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
		logger.Errorf(ctx, "Failed to register notification, %+v.", err)
		return nil, err
	}
	logger.Debugf(ctx, "Create notification [%s] in DB successfully.", notification.NotificationId)

	_, err = s.createTasksByNotification(ctx, notification)
	if err != nil {
		logger.Errorf(ctx, "Failed to create tasks by notification, %+v.", err)
		return nil, err
	}
	logger.Debugf(ctx, "Create tasks by notification [%s] in DB successfully.", notification.NotificationId)

	// Enqueue notification after create tasks.
	err = s.controller.notificationQueue.Enqueue(notification.NotificationId)
	if err != nil {
		logger.Errorf(ctx, "Push notification [%s] into etcd failed, %+v.", notification.NotificationId, err)
		return nil, err
	}
	logger.Debugf(ctx, "Push notification [%s] into etcd successfully.", notification.NotificationId)

	return &pb.CreateNotificationResponse{
		NotificationId: pbutil.ToProtoString(notification.NotificationId),
	}, nil
}

func (s *Server) createTasksByNotification(ctx context.Context, nf *models.Notification) ([]*models.Task, error) {
	tasks, err := rs.SplitNotificationIntoTasks(ctx, nf)
	if err != nil {
		logger.Errorf(ctx, "Failed to split notification into tasks, %+v.", err)
		return nil, err
	}

	err = s.createTasks(ctx, tasks)
	if err != nil {
		logger.Errorf(ctx, "Failed to create tasks, %+v.", err)
		return nil, err
	}

	return tasks, nil
}

func (s *Server) createTasks(ctx context.Context, tasks []*models.Task) error {
	var err error
	for _, task := range tasks {
		err = rs.RegisterTask(ctx, task)
		if err != nil {
			return err
		}
		logger.Debugf(ctx, "Create task [%s] in DB successfully.", task.TaskId)

		err = s.controller.taskQueue.Enqueue(task.TaskId)
		if err != nil {
			logger.Errorf(ctx, "Failed to push task [%s] into etcd, %+v.", task.TaskId, err)
			return err
		}
		logger.Debugf(ctx, "Push task [%s] into etcd successfully.", task.TaskId)
	}

	return nil
}

func (s *Server) RetryNotifications(ctx context.Context, req *pb.RetryNotificationsRequest) (*pb.RetryNotificationsResponse, error) {
	nfIds := stringutil.SimplifyStringList(req.NotificationId)
	nfs, err := rs.GetNfsByNfIds(ctx, req.NotificationId)
	if err != nil {
		logger.Errorf(ctx, "Failed to get notifications [%+v], %+v.", nfIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRetryTaskFailed, nfIds)
	}

	if len(nfs) != len(nfIds) {
		logger.Errorf(ctx, "Retry notifications [%+v] do not exit.", nfIds)
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorNotificationNotExist, nfIds)
	}

	for _, nf := range nfs {
		if !(nf.Status == constants.StatusSuccessful || nf.Status == constants.StatusFailed) {
			logger.Errorf(ctx, "Retry notifications [%+v] status is not %s or %s.", nfIds, constants.StatusSuccessful, constants.StatusFailed)
			return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorRetryNotificationsFailed, nfIds)
		}
	}

	for _, nfId := range nfIds {
		err = s.controller.notificationQueue.Enqueue(nfId)
		if err != nil {
			logger.Errorf(ctx, "Push notification [%s] into etcd failed, %+v.", nfId, err)
			return nil, err
		}
		logger.Debugf(ctx, "Push notification [%s] into etcd successfully.", nfId)
	}

	err = rs.UpdateNotificationsStatus(ctx, nfIds, constants.StatusSending)
	if err != nil {
		logger.Errorf(ctx, "Failed to update notifications [%+v] status to pending, %+v.", nfIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRetryNotificationsFailed, nfIds)
	}

	taskIds, err := rs.GetTaskIdsByNfIds(ctx, nfIds)
	if err != nil {
		logger.Errorf(ctx, "Failed to get task ids by notification ids[%+v], %+v.", nfIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRetryNotificationsFailed, nfIds)
	}

	err = s.retryTasksByTaskIds(ctx, taskIds)
	if err != nil {
		logger.Errorf(ctx, "Failed to retry tasks [%+v], %+v.", taskIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRetryNotificationsFailed, nfIds)
	}

	logger.Debugf(ctx, "Retry notifications [%+v] successfully.", nfIds)

	nfPbSet := models.NotificationSet2PbSet(nfs)
	res := &pb.RetryNotificationsResponse{
		NotificationSet: nfPbSet,
	}
	return res, nil
}

func (s *Server) RetryTasks(ctx context.Context, req *pb.RetryTasksRequest) (*pb.RetryTasksResponse, error) {
	taskIds := stringutil.SimplifyStringList(req.TaskId)
	tasks, err := rs.GetTasksByTaskIds(ctx, taskIds)
	if err != nil {
		logger.Errorf(ctx, "Failed to get tasks [%+v], %+v.", taskIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRetryTaskFailed, taskIds)
	}
	if len(tasks) != len(taskIds) {
		logger.Errorf(ctx, "Retry tasks [%+v] do not exit.", taskIds)
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorTaskNotExist, taskIds)
	}

	for _, task := range tasks {
		if !(task.Status == constants.StatusSuccessful || task.Status == constants.StatusFailed) {
			logger.Errorf(ctx, "Retry tasks [%+v] status is not %s or %s.", taskIds, constants.StatusSuccessful, constants.StatusFailed)
			return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorRetryTaskFailed, taskIds)
		}
	}

	err = s.retryTasksByTaskIds(ctx, taskIds)
	if err != nil {
		logger.Errorf(ctx, "Failed to retry tasks [%+v], %+v.", taskIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRetryTaskFailed, taskIds)
	}
	taskPbSet := models.TaskSet2PbSet(tasks)
	res := &pb.RetryTasksResponse{
		TaskSet: taskPbSet,
	}
	return res, nil
}
func (s *Server) retryTasksByTaskIds(ctx context.Context, taskIds []string) error {
	err := rs.UpdateTasksStatus(ctx, taskIds, constants.StatusPending)
	if err != nil {
		logger.Errorf(ctx, "Failed to update tasks [%+v] status to pending, %+v.", taskIds, err)
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRetryTaskFailed, taskIds)
	}
	for _, taskId := range taskIds {
		err = s.controller.taskQueue.Enqueue(taskId)
		if err != nil {
			logger.Errorf(ctx, "Failed to push task [%+v] into etcd, %+v.", taskIds, err)
			return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRetryTaskFailed, taskIds)
		}
		logger.Debugf(ctx, "Push task [%s] into etcd successfully.", taskId)
	}
	logger.Debugf(ctx, "Push tasks [%+v] into etcd successfully.", taskIds)
	return nil
}

func (s *Server) DescribeNotifications(ctx context.Context, req *pb.DescribeNotificationsRequest) (*pb.DescribeNotificationsResponse, error) {
	nfs, nfCnt, err := rs.DescribeNotifications(ctx, req)
	if err != nil {
		logger.Errorf(ctx, "Failed to describe notifications, %+v.", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	nfPbSet := models.NotificationSet2PbSet(nfs)
	res := &pb.DescribeNotificationsResponse{
		TotalCount:      uint32(nfCnt),
		NotificationSet: nfPbSet,
	}

	logger.Debugf(ctx, "Describe notifications successfully, notifications=%+v.", res)
	return res, nil
}

func (s *Server) DescribeTasks(ctx context.Context, req *pb.DescribeTasksRequest) (*pb.DescribeTasksResponse, error) {
	tasks, taskCnt, err := rs.DescribeTasks(ctx, req)
	if err != nil {
		logger.Errorf(ctx, "Failed to describe tasks,%+v.", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	taskPbSet := models.TaskSet2PbSet(tasks)
	res := &pb.DescribeTasksResponse{
		TotalCount: uint32(taskCnt),
		TaskSet:    taskPbSet,
	}
	logger.Debugf(ctx, "Describe tasks successfully, tasks=%+v.", res)
	return res, nil
}

func (s *Server) CreateAddress(ctx context.Context, req *pb.CreateAddressRequest) (*pb.CreateAddressResponse, error) {
	err := ValidateCreateAddressParams(ctx, req)
	if err != nil {
		return nil, err
	}

	addr := models.NewAddress(req)
	addrId, err := rs.CreateAddress(ctx, addr)
	if err != nil {
		logger.Errorf(ctx, "Failed to create address, %+v.", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	logger.Debugf(ctx, "Create address [%s] successfully.", addr.AddressId)

	return &pb.CreateAddressResponse{
		AddressId: pbutil.ToProtoString(addrId),
	}, nil
}

func (s *Server) DescribeAddresses(ctx context.Context, req *pb.DescribeAddressesRequest) (*pb.DescribeAddressesResponse, error) {
	addrs, addrCnt, err := rs.DescribeAddresses(ctx, req)
	if err != nil {
		logger.Errorf(ctx, "Failed to describe addresses, %+v.", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	addressPbSet := models.AddressSet2PbSet(addrs)
	res := &pb.DescribeAddressesResponse{
		TotalCount: uint32(addrCnt),
		AddressSet: addressPbSet,
	}
	logger.Debugf(ctx, "Describe addresses successfully, addrs = [%+v].", res)
	return res, nil
}

func (s *Server) ModifyAddress(ctx context.Context, req *pb.ModifyAddressRequest) (*pb.ModifyAddressResponse, error) {
	var err error
	err = ValidateModifyAddressParams(ctx, req)
	if err != nil {
		return nil, err
	}

	addressId := req.AddressId.GetValue()
	attributes := make(map[string]interface{})
	if req.Address.GetValue() != "" {
		attributes[models.AddrColAddress] = req.Address.GetValue()
	}
	if req.Remarks.GetValue() != "" {
		attributes[models.AddrColRemarks] = req.Remarks.GetValue()
	}
	if req.VerificationCode.GetValue() != "" {
		attributes[models.AddrColVerificationCode] = req.VerificationCode.GetValue()
	}
	if req.NotifyType.GetValue() != "" {
		attributes[models.AddrColNotifyType] = req.NotifyType.GetValue()
	}

	err = rs.ModifyAddress(ctx, addressId, attributes)
	if err != nil {
		logger.Errorf(ctx, "Failed to modify address [%s], %+v.", addressId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed, addressId)
	}

	logger.Debugf(ctx, "Modify address [%s] successfully.", addressId)

	return &pb.ModifyAddressResponse{
		AddressId: pbutil.ToProtoString(addressId),
	}, nil

}

func (s *Server) DeleteAddresses(ctx context.Context, req *pb.DeleteAddressesRequest) (*pb.DeleteAddressesResponse, error) {
	addressIds := stringutil.SimplifyStringList(req.AddressId)
	err := rs.DeleteAddresses(ctx, addressIds)
	if err != nil {
		logger.Errorf(ctx, "Failed to delete addresses [%+v], %+v.", addressIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourceFailed, addressIds)
	}

	logger.Debugf(ctx, "Delete addresses [%+v] successfully.", addressIds)
	return &pb.DeleteAddressesResponse{
		AddressId: addressIds,
	}, nil

}

func (s *Server) CreateAddressList(ctx context.Context, req *pb.CreateAddressListRequest) (*pb.CreateAddressListResponse, error) {
	addrList := models.NewAddressList(req)
	addrIds := req.GetAddressId()
	var addrListId string
	var err error

	if addrIds != nil {
		addrListId, err = rs.CreateAddressListWithAddrIDs(ctx, addrList, addrIds)
	} else {
		addrListId, err = rs.CreateAddressList(ctx, addrList)
	}
	if err != nil {
		logger.Errorf(ctx, "Failed to create addressList, %+v.", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	logger.Debugf(ctx, "Create addressList[%s] successfully.", addrList.AddressListId)

	return &pb.CreateAddressListResponse{
		AddressListId: pbutil.ToProtoString(addrListId),
	}, nil
}

func (s *Server) DescribeAddressList(ctx context.Context, req *pb.DescribeAddressListRequest) (*pb.DescribeAddressListResponse, error) {
	addrLists, addrCnt, err := rs.DescribeAddressLists(ctx, req)
	if err != nil {
		logger.Errorf(ctx, "Failed to describe addresses, %+v.", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	addressListPbSet := models.AddressListSet2PbSet(addrLists)
	res := &pb.DescribeAddressListResponse{
		TotalCount:     uint32(addrCnt),
		AddressListSet: addressListPbSet,
	}
	logger.Debugf(ctx, "Describe addressLists successfully, addrLists = [%+v].", res)
	return res, nil
}

func (s *Server) ModifyAddressList(ctx context.Context, req *pb.ModifyAddressListRequest) (*pb.ModifyAddressListResponse, error) {
	var err error
	addressListId := req.GetAddressListId().GetValue()
	addrIds := req.GetAddressId()

	attributes := make(map[string]interface{})
	if req.AddressListName.GetValue() != "" {
		attributes[models.AddrLsColName] = req.AddressListName.GetValue()
	}
	if req.Extra.GetValue() != "" {
		attributes[models.AddrLsColExtra] = req.Extra.GetValue()
	}
	if req.Status.GetValue() != "" {
		attributes[models.AddrLsColStatus] = req.Status.GetValue()
	}

	if addrIds != nil {
		err = rs.ModifyAddressListWithAddrIDs(ctx, addressListId, attributes, addrIds)
	} else {
		err = rs.ModifyAddressList(ctx, addressListId, attributes)
	}
	if err != nil {
		logger.Errorf(ctx, "Failed to modify addressList [%s], %+v.", addressListId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed, addressListId)
	}

	logger.Debugf(ctx, "Modify addressList [%s] successfully.", addressListId)

	return &pb.ModifyAddressListResponse{
		AddressListId: pbutil.ToProtoString(addressListId),
	}, nil
}

func (s *Server) DeleteAddressList(ctx context.Context, req *pb.DeleteAddressListRequest) (*pb.DeleteAddressListResponse, error) {
	addressListIds := stringutil.SimplifyStringList(req.AddressListId)
	err := rs.DeleteAddressLists(ctx, addressListIds)
	if err != nil {
		logger.Errorf(ctx, "Failed to delete address lists [%+v], %+v.", addressListIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourceFailed, addressListIds)
	}

	logger.Debugf(ctx, "Delete address list [%+v] successfully.", addressListIds)

	return &pb.DeleteAddressListResponse{
		AddressListId: addressListIds,
	}, nil

}

func (s *Server) ValidateEmailService(ctx context.Context, req *pb.ServiceConfig) (*pb.ValidateEmailServiceResponse, error) {
	host := req.GetEmailServiceConfig().GetEmailHost().GetValue()
	port := req.GetEmailServiceConfig().GetPort().GetValue()
	email := req.GetEmailServiceConfig().GetEmail().GetValue()
	password := req.GetEmailServiceConfig().GetPassword().GetValue()
	displaySender := req.GetEmailServiceConfig().GetDisplaySender().GetValue()

	emailAddr := email
	header := "ValidateEmailService"
	body := "<p>Email Service Config is working!</p>"

	m := gomail.NewMessage()
	m.SetAddressHeader("From", email, displaySender)
	m.SetHeader("To", emailAddr)
	m.SetHeader("Subject", header)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(host, int(port), email, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		logger.Errorf(ctx, "Send email to [%s] failed, [%+v]", emailAddr, err)
		return &pb.ValidateEmailServiceResponse{
			IsSucc: pbutil.ToProtoBool(false),
		}, err
	}

	return &pb.ValidateEmailServiceResponse{
		IsSucc: pbutil.ToProtoBool(true),
	}, nil

}
