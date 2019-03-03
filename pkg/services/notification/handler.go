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
	"openpitrix.io/notification/pkg/util/stringutil"
)

func (s *Server) SetServiceConfig(ctx context.Context, req *pb.ServiceConfig) (*pb.SetServiceConfigResponse, error) {
	err := ValidateSetServiceConfigParams(ctx, req)
	if err != nil {
		return nil, err
	}

	rs.SetServiceConfig(req)

	logger.Debugf(ctx, "Set ServiceConfig successfully, [%+v].", req)

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
		err := gerr.NewWithDetail(ctx, gerr.Internal, fmt.Errorf("Failed to Get ServiceConfig, EmailServiceConfig."), gerr.ErrorGetServiceConfigFailed)
		return nil, err
	}
	logger.Infof(ctx, "Get ServiceConfig successfully, [%+v].", emailCfg)

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
	logger.Debugf(ctx, "Create notification[%s] in DB successfully.", notification.NotificationId)

	_, err = s.createTasksByNotification(ctx, notification)
	if err != nil {
		logger.Errorf(ctx, "Create tasks by notification failed, [%+v].", err)
		return nil, err
	}

	// Enqueue notification after create tasks.
	err = s.controller.notificationQueue.Enqueue(notification.NotificationId)
	if err != nil {
		logger.Errorf(ctx, "Push notification[%s] into etcd failed, [%+v].", notification.NotificationId, err)
		return nil, err
	}
	logger.Debugf(ctx, "Push notification[%s] into etcd successfully.", notification.NotificationId)
	logger.Debugf(ctx, "Create notification[%s] successfully.", notification.NotificationId)

	return &pb.CreateNotificationResponse{
		NotificationId: pbutil.ToProtoString(notification.NotificationId),
	}, nil
}

func (s *Server) createTasksByNotification(ctx context.Context, nf *models.Notification) ([]*models.Task, error) {
	tasks, err := rs.SplitNotificationIntoTasks(nf)
	if err != nil {
		logger.Errorf(ctx, "Split notification into tasks failed, [%+v].", err)
		return nil, err
	}

	tasks, err = s.createTasks(ctx, tasks)
	if err != nil {
		logger.Errorf(ctx, "Split notification into tasks failed, [%+v].", err)
		return nil, err
	}
	return tasks, nil
}

func (s *Server) createTasks(ctx context.Context, tasks []*models.Task) ([]*models.Task, error) {
	var err error
	for _, task := range tasks {
		err = rs.RegisterTask(ctx, task)
		if err != nil {
			return nil, err
		}
		logger.Debugf(ctx, "Create Task[%s] in DB successfully.", task.TaskId)

		err = s.controller.taskQueue.Enqueue(task.TaskId)
		if err != nil {
			logger.Errorf(ctx, "Push task[%s] into etcd failed, [%+v].", task.TaskId, err)
			return nil, err
		}
		logger.Debugf(ctx, "Push task[%s] into etcd successfully.", task.TaskId)
		logger.Debugf(ctx, "Create Task[%s] successfully.", task.TaskId)
	}

	return tasks, nil
}

func (s *Server) DescribeNotifications(ctx context.Context, req *pb.DescribeNotificationsRequest) (*pb.DescribeNotificationsResponse, error) {
	nfs, nfCnt, err := rs.DescribeNotifications(ctx, req)
	if err != nil {
		logger.Errorf(ctx, "Failed to describe notifications, [%+v].", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	nfPbSet := models.ParseNfSet2PbSet(nfs)
	res := &pb.DescribeNotificationsResponse{
		TotalCount:      uint32(nfCnt),
		NotificationSet: nfPbSet,
	}

	logger.Debugf(ctx, "Describe notifications successfully, notifications=[%+v].", res)
	return res, nil
}

func (s *Server) DescribeTasks(ctx context.Context, req *pb.DescribeTasksRequest) (*pb.DescribeTasksResponse, error) {
	tasks, taskCnt, err := rs.DescribeTasks(ctx, req)
	if err != nil {
		logger.Errorf(ctx, "Failed to describe tasks,[%+v].", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	taskPbSet := models.ParseTaskSet2PbSet(tasks)
	res := &pb.DescribeTasksResponse{
		TotalCount: uint32(taskCnt),
		TaskSet:    taskPbSet,
	}
	logger.Debugf(ctx, "Describe tasks successfully, tasks=[%+v].", res)
	return res, nil
}

func (s *Server) RetryNotifications(ctx context.Context, req *pb.RetryNotificationsRequest) (*pb.RetryNotificationsResponse, error) {
	nfIds := stringutil.SimplifyStringList(req.NotificationId)
	nfs, err := rs.GetNfsByNfIds(req.NotificationId)
	if err != nil {
		logger.Errorf(ctx, "Failed to retry notifications[%+v], [%+v].", nfIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRetryTaskFailed, nfIds)
	}
	if len(nfs) == 0 {
		logger.Infof(ctx, "Retry notifications[%+v] do not exit.", nfIds)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRetryNotificationtNotExist, nfIds)
	}

	nfIds, err = rs.UpdateNotifications2Pending(ctx, stringutil.SimplifyStringList(req.NotificationId))
	if err != nil {
		logger.Errorf(ctx, "Failed to retry notifications[%+v], [%+v].", nfIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRetryNotificationsFailed, nfIds)
	}

	tasks, err := rs.GetTasksByNfId(nfIds)
	var taskIds []string
	for _, task := range tasks {
		taskId := task.TaskId
		taskIds = append(taskIds, taskId)
	}

	_, err = s.retryTasksByTaskIds(ctx, taskIds)
	if err != nil {
		logger.Errorf(ctx, "Failed to retry notifications[%+v], [%+v].", nfIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRetryNotificationsFailed, nfIds)
	}

	var nfReq = &pb.DescribeNotificationsRequest{
		NotificationId: nfIds,
	}
	nfs, _, err = rs.DescribeNotifications(ctx, nfReq)
	if err != nil {
		logger.Errorf(ctx, "Failed to retry notifications[%+v], [%+v].", nfIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRetryNotificationsFailed, nfIds)
	}

	logger.Debugf(ctx, "Retry notifications[%+v] successfully.", nfIds)

	nfPbSet := models.ParseNfSet2PbSet(nfs)
	res := &pb.RetryNotificationsResponse{
		NotificationSet: nfPbSet,
	}
	return res, nil
}

func (s *Server) RetryTasks(ctx context.Context, req *pb.RetryTasksRequest) (*pb.RetryTasksResponse, error) {
	taskIds := stringutil.SimplifyStringList(req.TaskId)
	tasks, err := rs.GetTasksByTaskIds(taskIds)
	if err != nil {
		logger.Errorf(ctx, "Failed to retry tasks[%+v], [%+v].", taskIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRetryTaskFailed, taskIds)
	}
	if len(tasks) == 0 {
		logger.Infof(ctx, "Retry tasks[%+v] do not exit.", taskIds)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRetryTaskNotExist, taskIds)
	}

	tasks, err = s.retryTasksByTaskIds(ctx, taskIds)
	if err != nil {
		logger.Errorf(ctx, "Failed to retry tasks[%+v], [%+v].", taskIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRetryTaskFailed, taskIds)
	}
	taskPbSet := models.ParseTaskSet2PbSet(tasks)
	res := &pb.RetryTasksResponse{
		TaskSet: taskPbSet,
	}
	return res, nil
}
func (s *Server) retryTasksByTaskIds(ctx context.Context, taskIds []string) ([]*models.Task, error) {
	taskIds, err := rs.UpdateTasks2Pending(ctx, taskIds)
	if err != nil {
		logger.Errorf(ctx, "Failed to retry tasks[%+v], [%+v].", taskIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRetryTaskFailed, taskIds)
	}

	var taskReq = &pb.DescribeTasksRequest{
		TaskId: taskIds,
	}
	tasks, _, err := rs.DescribeTasks(ctx, taskReq)
	if err != nil {
		logger.Errorf(ctx, "Failed to retry tasks[%+v], [%+v].", taskIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRetryTaskFailed, taskIds)
	}

	for _, task := range tasks {
		err = s.controller.taskQueue.Enqueue(task.TaskId)
		if err != nil {
			logger.Errorf(ctx, "Failed to retry tasks[%+v], [%+v].", taskIds, err)
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRetryTaskFailed, taskIds)
		}
		logger.Debugf(ctx, "Push task[%s] into etcd successfully.", task.TaskId)
	}
	return tasks, nil
}

func (s *Server) CreateAddress(ctx context.Context, req *pb.CreateAddressRequest) (*pb.CreateAddressResponse, error) {
	err := ValidateCreateAddressParams(ctx, req)
	if err != nil {
		return nil, err
	}

	addr := models.NewAddress(req)
	addrId, err := rs.CreateAddress(ctx, addr)

	if err != nil {
		logger.Errorf(ctx, "Failed to create address[%+v], [%+v].", addr, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}
	logger.Debugf(ctx, "Create address[%s] successfully.", addr.AddressId)
	return &pb.CreateAddressResponse{
		AddressId: pbutil.ToProtoString(addrId),
	}, nil
}

func (s *Server) DescribeAddresses(ctx context.Context, req *pb.DescribeAddressesRequest) (*pb.DescribeAddressesResponse, error) {
	addrs, addrCnt, err := rs.DescribeAddresses(ctx, req)
	if err != nil {
		logger.Errorf(ctx, "Failed to describe addresses,[%+v].", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	addressPbSet := models.ParseAddressSet2PbSet(addrs)
	res := &pb.DescribeAddressesResponse{
		TotalCount: uint32(addrCnt),
		AddressSet: addressPbSet,
	}
	logger.Debugf(ctx, "Describe addresses successfully, addrs=[%+v].", res)
	return res, nil
}

func (s *Server) ModifyAddress(ctx context.Context, req *pb.ModifyAddressRequest) (*pb.ModifyAddressResponse, error) {
	err := ValidateModifyAddressParams(ctx, req)
	if err != nil {
		return nil, err
	}

	addressId, err := rs.ModifyAddress(ctx, req)
	if err != nil {
		logger.Errorf(ctx, "Failed to modify address[%s], [%+v].", addressId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed, addressId)
	}
	logger.Debugf(ctx, "Modify address[%s] successfully.", addressId)
	return &pb.ModifyAddressResponse{
		AddressId: pbutil.ToProtoString(addressId),
	}, nil

}

func (s *Server) DeleteAddresses(ctx context.Context, req *pb.DeleteAddressesRequest) (*pb.DeleteAddressesResponse, error) {
	addressIds, err := rs.DeleteAddresses(ctx, stringutil.SimplifyStringList(req.AddressId))
	if err != nil {
		logger.Errorf(ctx, "Failed to delete addresses[%+v], [%+v].", addressIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourceFailed, addressIds)
	}
	logger.Debugf(ctx, "Delete addresses[%+v] successfully.", addressIds)
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
		logger.Errorf(ctx, "Failed to create addressList[%+v].", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	logger.Debugf(ctx, "Create addressList[%+v] successfully.", addrList)

	return &pb.CreateAddressListResponse{
		AddressListId: pbutil.ToProtoString(addrListId),
	}, nil
}

func (s *Server) DescribeAddressList(ctx context.Context, req *pb.DescribeAddressListRequest) (*pb.DescribeAddressListResponse, error) {
	addrLists, addrCnt, err := rs.DescribeAddressLists(ctx, req)
	if err != nil {
		logger.Errorf(ctx, "Failed to describe addresses,[%+v].", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	addressListPbSet := models.ParseAddressListSet2PbSet(addrLists)
	res := &pb.DescribeAddressListResponse{
		TotalCount:     uint32(addrCnt),
		AddressListSet: addressListPbSet,
	}
	logger.Debugf(ctx, "Describe addressLists successfully, addrLists=[%+v].", res)
	return res, nil
}

func (s *Server) ModifyAddressList(ctx context.Context, req *pb.ModifyAddressListRequest) (*pb.ModifyAddressListResponse, error) {
	addressListId, err := rs.ModifyAddressList(ctx, req)
	if err != nil {
		logger.Errorf(ctx, "Failed to modify addressList[%s], [%+v].", addressListId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed, addressListId)
	}
	logger.Debugf(ctx, "Modify addressList[%s] successfully.", addressListId)
	return &pb.ModifyAddressListResponse{
		AddressListId: pbutil.ToProtoString(addressListId),
	}, nil
}

func (s *Server) DeleteAddressList(ctx context.Context, req *pb.DeleteAddressListRequest) (*pb.DeleteAddressListResponse, error) {
	addressListIds, err := rs.DeleteAddressLists(ctx, stringutil.SimplifyStringList(req.AddressListId))
	if err != nil {
		logger.Errorf(ctx, "Failed to delete address Lists[%+v], [%+v].", addressListIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourceFailed, addressListIds)
	}
	logger.Debugf(ctx, "Delete addressList[%+v] successfully.", addressListIds)
	return &pb.DeleteAddressListResponse{
		AddressListId: addressListIds,
	}, nil

}
