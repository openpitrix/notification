// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"context"
	"fmt"

	"openpitrix.io/logger"

	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/gerr"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/pb"
	rs "openpitrix.io/notification/pkg/services/notification/resource_control"
	"openpitrix.io/notification/pkg/util/emailutil"
	"openpitrix.io/notification/pkg/util/pbutil"
	"openpitrix.io/notification/pkg/util/stringutil"
)

//****************************************
//ServiceConfig
//****************************************
func (s *Server) SetServiceConfig(ctx context.Context, req *pb.ServiceConfig) (*pb.SetServiceConfigResponse, error) {
	err := ValidateSetServiceConfigParams(ctx, req.GetEmailServiceConfig(), "")
	if err != nil {
		logger.Errorf(ctx, "Failed to set service config.")
		return nil, gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorValidateEmailService)
	}

	err = rs.ModifyEmailConfig(ctx, req)
	if err != nil {
		logger.Errorf(ctx, "Failed to set service config, %+v.", err)
		return nil, err
	}

	logger.Debugf(ctx, "Set service config successfully, %+v.", config.GetInstance().Email)
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

	var emailCfgPb *pb.EmailServiceConfig
	scCfg := &pb.ServiceConfig{}
	for _, scType := range serviceTypes {
		if scType == constants.NotifyTypeEmail {
			emailCfg, err := rs.GetEmailConfigFromDB(ctx)
			if err != nil {
				logger.Errorf(ctx, "Failed to get email config, %+v.", err)
				return nil, err
			}
			emailCfgPb = models.EmailConfigToPb(emailCfg)
			logger.Debugf(ctx, "Get service config [%+v] successfully.", emailCfg)
			break
		}
	}
	if emailCfgPb == nil {
		logger.Errorf(ctx, "Failed to get service config, email service config.")
		err := gerr.NewWithDetail(ctx, gerr.Internal, fmt.Errorf("Failed to get service config, email service config."), gerr.ErrorGetServiceConfigFailed)
		return nil, err
	}
	logger.Debugf(ctx, "Get service config [%+v] successfully.", emailCfgPb)

	scCfg.EmailServiceConfig = emailCfgPb
	return scCfg, nil
}

func (s *Server) ValidateEmailService(ctx context.Context, req *pb.ServiceConfig) (*pb.ValidateEmailServiceResponse, error) {
	err := ValidateSetServiceConfigParams(ctx, req.GetEmailServiceConfig(), "")
	if err != nil {
		return nil, err
	}

	err = emailutil.SendMail4ValidateEmailService(nil, req.GetEmailServiceConfig(), "", "zh")
	if err != nil {
		logger.Errorf(nil, "send email failed, [%+v]", err)
		return &pb.ValidateEmailServiceResponse{
			IsSucc: pbutil.ToProtoBool(false),
		}, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorValidateEmailService)
	}
	return &pb.ValidateEmailServiceResponse{
		IsSucc: pbutil.ToProtoBool(true),
	}, nil
}

func (s *Server) ValidateEmailServiceV2(ctx context.Context, req *pb.ValidateEmailServiceV2Request) (*pb.ValidateEmailServiceResponse, error) {
	err := ValidateSetServiceConfigParams(ctx, req.GetEmailServiceConfig(), req.GetTestEmailRecipient().GetValue())
	if err != nil {
		return nil, err
	}

	language := req.GetLanguage().GetValue()
	if language == "" {
		language = "zh"
	}

	err = emailutil.SendMail4ValidateEmailService(nil, req.GetEmailServiceConfig(), req.GetTestEmailRecipient().GetValue(), language)
	if err != nil {
		logger.Errorf(nil, "send email failed, [%+v]", err)
		return &pb.ValidateEmailServiceResponse{
			IsSucc: pbutil.ToProtoBool(false),
		}, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorValidateEmailService)
	}
	return &pb.ValidateEmailServiceResponse{
		IsSucc: pbutil.ToProtoBool(true),
	}, nil
}

//****************************************
//Notification
//****************************************
func (s *Server) CreateNotification(ctx context.Context, req *pb.CreateNotificationRequest) (*pb.CreateNotificationResponse, error) {
	//Step0:Validate params.
	err := ValidateCreateNotification(ctx, req)
	if err != nil {
		return nil, err
	}

	notification := models.NewNotification(req)

	//Step1:Register Notification in DB as status="pending"
	err = rs.RegisterNotification(ctx, notification)
	if err != nil {
		logger.Errorf(ctx, "Failed to register notification, %+v.", err)
		return nil, err
	}
	logger.Debugf(ctx, "Create notification[%s] in DB successfully.", notification.NotificationId)

	//Step2:Process Task data.
	//including:
	// 2.1.SplitNotificationIntoTasks,
	// 2.2.createTasks in db,
	// 2.3.Enqueue task id to queue.
	// 2.4 if it is websocket messge, publish to pubsub.
	_, err = s.createTasksByNotification(ctx, notification)
	if err != nil {
		logger.Errorf(ctx, "Failed to create tasks by notification, %+v.", err)
		return nil, err
	}
	logger.Debugf(ctx, "Create tasks by notification[%s] in DB successfully.", notification.NotificationId)

	//Step3:Enqueue notification id to queue.
	err = s.controller.notificationQueue.Enqueue(notification.NotificationId)

	if err != nil {
		logger.Errorf(ctx, "Push notification[%s] into queue failed, %+v.", notification.NotificationId, err)
		return nil, err
	}
	logger.Debugf(ctx, "Push notification[%s] into queue successfully.", notification.NotificationId)

	return &pb.CreateNotificationResponse{
		NotificationId: pbutil.ToProtoString(notification.NotificationId),
	}, nil
}
func (s *Server) createTasksByNotification(ctx context.Context, nf *models.Notification) ([]*models.Task, error) {
	//1.SplitNotificationIntoTasks
	//2.createTasks
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
		logger.Debugf(ctx, "Create task[%s] in DB successfully.", task.TaskId)

		if task.NotifyType == constants.NotifyTypeEmail {
			err = s.controller.taskQueue.Enqueue(task.TaskId)
			if err != nil {
				logger.Errorf(ctx, "Failed to push task[%s] into queue, %+v.", task.TaskId, err)
				return err
			}
			logger.Debugf(ctx, "Push task[%s] into queue successfully.", task.TaskId)
		}
	}

	return nil
}

func (s *Server) RetryNotifications(ctx context.Context, req *pb.RetryNotificationsRequest) (*pb.RetryNotificationsResponse, error) {
	nfIds := stringutil.SimplifyStringList(req.NotificationId)
	nfs, err := rs.GetFailedNfsByNfIds(ctx, req.NotificationId)
	if err != nil {
		logger.Errorf(ctx, "Failed to get notifications[%+v], %+v.", nfIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRetryTaskFailed, nfIds)
	}

	if len(nfs) != len(nfIds) {
		logger.Errorf(ctx, "Retry notifications[%+v] do not exit.", nfIds)
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotExist, nfIds)
	}

	for _, nf := range nfs {
		if !(nf.Status == constants.StatusSuccessful || nf.Status == constants.StatusFailed) {
			logger.Errorf(ctx, "Retry notifications[%+v] status is not %s or %s.", nfIds, constants.StatusSuccessful, constants.StatusFailed)
			return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorRetryNotificationsFailed, nfIds)
		}
	}

	for _, nfId := range nfIds {
		err = s.controller.notificationQueue.Enqueue(nfId)
		if err != nil {
			logger.Errorf(ctx, "Push notification[%s] into queue failed, %+v.", nfId, err)
			return nil, err
		}
		logger.Debugf(ctx, "Push notification[%s] into queue successfully.", nfId)
	}

	err = rs.UpdateNotificationsStatus(ctx, nfIds, constants.StatusSending)
	if err != nil {
		logger.Errorf(ctx, "Failed to update notifications[%+v] status to pending, %+v.", nfIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRetryNotificationsFailed, nfIds)
	}

	taskIds, err := rs.GetTaskIdsByNfIds(ctx, nfIds)
	if err != nil {
		logger.Errorf(ctx, "Failed to get task ids by notification ids[%+v], %+v.", nfIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRetryNotificationsFailed, nfIds)
	}

	err = s.retryTasksByTaskIds(ctx, taskIds)
	if err != nil {
		logger.Errorf(ctx, "Failed to retry tasks[%+v], %+v.", taskIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRetryNotificationsFailed, nfIds)
	}

	logger.Debugf(ctx, "Retry notifications[%+v] successfully.", nfIds)

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
		logger.Errorf(ctx, "Failed to get tasks[%+v], %+v.", taskIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRetryTaskFailed, taskIds)
	}
	if len(tasks) != len(taskIds) {
		logger.Errorf(ctx, "Retry tasks[%+v] do not exit.", taskIds)
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotExist, taskIds)
	}

	for _, task := range tasks {
		if !(task.Status == constants.StatusSuccessful || task.Status == constants.StatusFailed) {
			logger.Errorf(ctx, "Retry tasks[%+v] status is not %s or %s.", taskIds, constants.StatusSuccessful, constants.StatusFailed)
			return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorRetryTaskFailed, taskIds)
		}
	}

	err = s.retryTasksByTaskIds(ctx, taskIds)
	if err != nil {
		logger.Errorf(ctx, "Failed to retry tasks[%+v], %+v.", taskIds, err)
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
		logger.Errorf(ctx, "Failed to update tasks[%+v] status to pending, %+v.", taskIds, err)
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRetryTaskFailed, taskIds)
	}
	for _, taskId := range taskIds {
		err = s.controller.taskQueue.Enqueue(taskId)
		if err != nil {
			logger.Errorf(ctx, "Failed to push task[%+v] into queue, %+v.", taskIds, err)
			return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRetryTaskFailed, taskIds)
		}
		logger.Debugf(ctx, "Push task[%s] into queue successfully.", taskId)
	}
	logger.Debugf(ctx, "Push tasks[%+v] into queue successfully.", taskIds)
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

	logger.Debugf(ctx, "Describe notifications successfully, notifications count=%+v.", res.TotalCount)
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
	logger.Debugf(ctx, "Describe tasks successfully, tasks count=%+v.", res.TotalCount)

	return res, nil
}

//****************************************
//Address
//****************************************
func (s *Server) CreateAddress(ctx context.Context, req *pb.CreateAddressRequest) (*pb.CreateAddressResponse, error) {
	//step1:validate
	err := ValidateCreateAddress(ctx, req)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorValidateFailed)
	}
	//Step2:if address does not exist,create it.
	resp, err := s.createAddress(ctx, req)
	if err != nil {
		logger.Errorf(ctx, "Failed to create address, %+v.", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}
	return resp, nil
}
func (s *Server) createAddress(ctx context.Context, req *pb.CreateAddressRequest) (*pb.CreateAddressResponse, error) {
	addr := models.NewAddress(req)
	addrId, err := rs.CreateAddress(ctx, addr)
	if err != nil {
		logger.Errorf(ctx, "Failed to create address, %+v.", err)
		return nil, err
	}
	logger.Debugf(ctx, "Create address[%s] successfully.", addr.AddressId)
	return &pb.CreateAddressResponse{
		AddressId: pbutil.ToProtoString(addrId),
	}, nil
}

func (s *Server) DescribeAddresses(ctx context.Context, req *pb.DescribeAddressesRequest) (*pb.DescribeAddressesResponse, error) {
	addrs, addrCnt, err := rs.DescribeAddressesWithAddrListId(ctx, req)
	if err != nil {
		logger.Errorf(ctx, "Failed to describe addresses, %+v.", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	addressPbSet := models.AddressWithListIdSet2PbSet(addrs)
	res := &pb.DescribeAddressesResponse{
		TotalCount: uint32(addrCnt),
		AddressSet: addressPbSet,
	}
	logger.Debugf(ctx, "Describe addresses successfully, addrs = [%+v].", res)
	return res, nil

}

func (s *Server) ModifyAddress(ctx context.Context, req *pb.ModifyAddressRequest) (*pb.ModifyAddressResponse, error) {
	var err error
	//step1:validate
	err = ValidateModifyAddress(ctx, req)
	if err != nil {
		return nil, err
	}

	//step2:update, at first check record exist or not,
	//if not exist throw error
	//if exist lock the records in tx, then update.
	addressId := req.Address
	attributes := make(map[string]interface{})
	if req.AddressDetail.GetValue() != "" {
		attributes[models.AddrColAddress] = req.AddressDetail.GetValue()
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
		logger.Errorf(ctx, "Failed to modify address[%s], %+v.", addressId, err)
		return nil, err
	}
	logger.Debugf(ctx, "Modify address[%s] successfully.", addressId)

	return &pb.ModifyAddressResponse{
		AddressId: pbutil.ToProtoString(addressId),
	}, nil

}

func (s *Server) DeleteAddresses(ctx context.Context, req *pb.DeleteAddressesRequest) (*pb.DeleteAddressesResponse, error) {
	//step1:validate
	err := ValidateDeleteAddresses(ctx, req)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorValidateFailed)
	}

	//step2:delete
	addrIds := stringutil.Unique(stringutil.SimplifyStringList(req.GetAddressId()))
	err = rs.DeleteAddresses(ctx, addrIds)
	if err != nil {
		logger.Errorf(ctx, "Failed to delete addresses[%+v], %+v.", addrIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourceFailed, addrIds)
	}

	logger.Debugf(ctx, "Delete addresses[%+v] successfully.", addrIds)
	return &pb.DeleteAddressesResponse{
		AddressId: addrIds,
	}, nil

}

//****************************************
//AddressList
//****************************************
func (s *Server) CreateAddressList(ctx context.Context, req *pb.CreateAddressListRequest) (*pb.CreateAddressListResponse, error) {
	//step1:validate
	var err error
	err = ValidateCreateAddressList(ctx, req)
	if err != nil {
		return nil, err
	}

	//step2:create
	addrList := models.NewAddressList(req)
	addrIds := stringutil.Unique(stringutil.SimplifyStringList(req.GetAddressId()))
	var addrListId string

	if addrIds != nil {
		addrListId, err = rs.CreateAddressListWithAddrIDs(ctx, addrList, addrIds)
	} else {
		addrListId, err = rs.CreateAddressList(ctx, addrList)
	}
	if err != nil {
		logger.Errorf(ctx, "Failed to create addressList, %+v.", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	logger.Debugf(ctx, "Create addressList[%s] successfully.", addrListId)
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
	addressListId := req.GetAddresslist()
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
		err = ValidatModifyAddressList(ctx, req)
		if err != nil {
			return nil, err
		}
		err = rs.ModifyAddressListWithAddrIDs(ctx, addressListId, attributes, addrIds)
	} else {
		err = rs.ModifyAddressList(ctx, addressListId, attributes)
	}
	if err != nil {
		logger.Errorf(ctx, "Failed to modify addressList[%s], %+v.", addressListId, err)
		return nil, err
	}

	logger.Debugf(ctx, "Modify addressList[%s] successfully.", addressListId)

	return &pb.ModifyAddressListResponse{
		AddressListId: pbutil.ToProtoString(addressListId),
	}, nil
}

func (s *Server) DeleteAddressList(ctx context.Context, req *pb.DeleteAddressListRequest) (*pb.DeleteAddressListResponse, error) {
	//step1:validate
	err := ValidateDeleteAddressesList(ctx, req)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorValidateFailed)
	}
	//step2:delete
	addressListIds := stringutil.Unique(stringutil.SimplifyStringList(req.AddressListId))
	err = rs.DeleteAddressLists(ctx, addressListIds)
	if err != nil {
		logger.Errorf(ctx, "Failed to delete address lists[%+v], %+v.", addressListIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourceFailed, addressListIds)
	}

	logger.Debugf(ctx, "Delete address list[%+v] successfully.", addressListIds)

	return &pb.DeleteAddressListResponse{
		AddressListId: addressListIds,
	}, nil

}

func (s *Server) CreateNotificationChannel(req *pb.StreamReqData, res pb.Notification_CreateNotificationChannelServer) error {
	service := req.GetService().GetValue()

	for {
		outMsg := <-s.controller.websocketMsgChanMap[service]
		userMsg, err := models.UseMsgStringToPb(outMsg)
		if err != nil {
			logger.Errorf(nil, "Decode user message string to pb failed,err=%+v", err)
		}
		if service == userMsg.GetService().GetValue() {
			streamRespData := pb.StreamRespData{
				UserMsg: userMsg,
			}
			res.Send(&streamRespData)
		}

	}
}
