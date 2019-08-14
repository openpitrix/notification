// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"context"
	"regexp"
	"strconv"
	"time"

	"openpitrix.io/logger"

	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/gerr"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/pb"
	rs "openpitrix.io/notification/pkg/services/notification/resource_control"
	"openpitrix.io/notification/pkg/util/stringutil"
)

//****************************************
//Validate ServiceConfig
//****************************************
func ValidateSetServiceConfigParams(ctx context.Context, req *pb.ServiceConfig) error {
	email := req.GetEmailServiceConfig().GetEmail().GetValue()
	err := VerifyEmailFmt(ctx, email)
	if err != nil {
		logger.Errorf(ctx, "Failed to validate email [%s]: %+v", email, err)
		return err
	}

	port := req.GetEmailServiceConfig().GetPort().GetValue()
	err = VerifyPortFmt(ctx, int32(port))
	if err != nil {
		logger.Errorf(ctx, "Failed to validate port [%d]: %+v", port, err)
		return err
	}
	return nil
}

//****************************************
//Validate Notification
//****************************************
func ValidateCreateNotification(ctx context.Context, req *pb.CreateNotificationRequest) error {
	//1.validate avaiblable time
	err := validateAvaiblableTime(ctx, req)
	if err != nil {
		return err
	}

	//2.validate Content fmt
	//{\"html\":\"test_content_html\",  \"normal\":\"test_content_normal\"}
	_, err = models.DecodeContent(req.GetContent().GetValue())
	if err != nil {
		return err
	}
	notification := models.NewNotification(req)

	//3. validate address info
	_, decodeMapErr := models.DecodeAddressInfo(notification.AddressInfo)
	if decodeMapErr == nil {
		//3.1. check addressInfo format is like address_info = {"email": ["xxx@abc.com", "xxx@xxx.com"],"websocket": ["system", "huojiao"]}
		//if needed, also check extra fmt eg:"{"ws_service": "ks","ws_message_type": "event"}"
		err = validateAddressInfo4AddressMap(ctx, notification)
		if err != nil {
			return err
		}
	} else {
		//3.2 check addressInfo format is like address_info = ["adl-xxxx1", "adl-xxxx2"]
		//if address_info = ["adl-xxxx1", "adl-xxxx2"],check adl exists or not
		//todo if needed, also check extra fmt eg:"{"ws_service": "ks","ws_message_type": "event"}"
		err = validateAddressInfo4AddressListIds(ctx, notification)
		if err != nil {
			return err
		}
	}

	return nil
}
func validateAvaiblableTime(ctx context.Context, req *pb.CreateNotificationRequest) error {
	if req.GetAvailableStartTime().GetValue() != "" {
		availableStartTimeStr := req.GetAvailableStartTime().GetValue()
		err := VerifyAvailableTimeStr(ctx, availableStartTimeStr)
		if err != nil {
			logger.Errorf(ctx, "Failed to validate available start time [%s]: %+v", availableStartTimeStr, err)
			return err
		}
	}
	if req.GetAvailableStartTime().GetValue() != "" {
		availableEndTimeStr := req.AvailableEndTime.GetValue()
		err := VerifyAvailableTimeStr(ctx, availableEndTimeStr)
		if err != nil {
			logger.Errorf(ctx, "Failed to validate available end time [%s]: %+v", availableEndTimeStr, err)
			return err
		}
	}
	return nil
}
func validateAddressInfo4AddressListIds(ctx context.Context, notification *models.Notification) error {
	addressListIds, err := models.DecodeAddressListIds(notification.AddressInfo)
	if err == nil {
		//if address_info = ["adl-xxxx1", "adl-xxxx2"],check adl ids exists in DB or not
		*addressListIds = stringutil.Unique(*addressListIds)
		err = checkAddrListExistInDB(ctx, *addressListIds)
		if err != nil {
			logger.Errorf(ctx, "Failed to validate addressInfo, some address lists in address info[%s] do not exits: %+v", notification.AddressInfo, err)
			return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorNotExistItemInList, notification.AddressInfo)
		}

		if len(*addressListIds) == 0 {
			logger.Errorf(ctx, "Failed to validate addressInfo, address list id is blank: %+v", err)
			return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalNotificationAddressList, notification.AddressInfo)
		}
		for _, addressListId := range *addressListIds {
			//check addresslist id start with adl-, address_info = ["adl-xxxx1", "adl-xxxx2"]
			prefix := string([]byte(addressListId)[:4])
			if prefix != models.AddressListIdPrefix {
				logger.Errorf(ctx, "Failed to validate addressInfo[%s], address list id is should be [\"adl-xxxx1\", \"adl-xxxx2\"] ", addressListId)
				return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalNotificationAddressList, notification.AddressInfo)
			}
		}

	} else {
		logger.Errorf(ctx, "Failed to validate addressInfo[%s]: %+v", notification.AddressInfo, err)
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalNotificationAddressInfo)
	}
	return nil
}
func checkAddrListExistInDB(ctx context.Context, addressListIds []string) error {
	addrLists, err := rs.GetActiveAddressesListsByIds(ctx, addressListIds)
	if len(addrLists) != len(addressListIds) {
		return err
	}
	return nil
}
func validateAddressInfo4AddressMap(ctx context.Context, notification *models.Notification) error {
	//check addressinfo could be decode to map[string][]string
	addrInfoMap, err := models.DecodeAddressInfo(notification.AddressInfo)
	if err == nil {
		if len(*addrInfoMap) == 0 {
			logger.Errorf(ctx, "Failed to validate addressInfo,address info is blank: [%s].", notification.AddressInfo)
			return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalNotificationAddressInfo)
		}

		//check  address_info does not include websocket, throw error to make sure extra is not set or is {}
		_, ok := (*addrInfoMap)[constants.NotifyTypeWebsocket]
		if !ok {
			if notification.Extra != "{}" {
				return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalNotificationExtraBlank)
			}
		}

		//check addressInfo format is like address_info = {"email": ["xxx@abc.com", "xxx@xxx.com"],"websocket": ["system", "huojiao"]}
		for notifyType, values := range *addrInfoMap {
			if !(notifyType == constants.NotifyTypeEmail || notifyType == constants.NotifyTypeWebsocket) {
				logger.Errorf(ctx, "Failed to validate addressInfo, notify type is invalid [%s].", notifyType)
				return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalNotificationType, notifyType)
			}
			if notifyType == constants.NotifyTypeEmail {
				for _, email := range values {
					err := VerifyEmailFmt(ctx, email)
					if err != nil {
						logger.Errorf(ctx, "Failed to validate email [%s]: %+v", email, err)
						return err
					}
				}
			}
			if notifyType == constants.NotifyTypeWebsocket {
				//At first shuld check the websocket config is enabled
				if config.GetInstance().Websocket.Service == "none" {
					logger.Errorf(ctx, "Failed to validate addressInfo,websocket config is disabled.")
					return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalWebsocketDisabled, notification.AddressInfo)
				}
				//extra info is only used for websocket notification,to show which websocket client could accept it.
				//eg:"{"ws_service": "ks","ws_message_type": "event"}"
				err := models.CheckExtra(ctx, notification.Extra)
				if err != nil {
					return gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorDecodeExtraFailed)
				}
			}
		}
		return nil

	} else {
		return gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorIllegalNotificationAddressInfo)
	}

}

//****************************************
//Validate AddressList
//****************************************
func ValidateCreateAddressList(ctx context.Context, req *pb.CreateAddressListRequest) error {
	addrIds := stringutil.Unique(req.GetAddressId())
	//check addr status is active,only active addr could be add to a new address list.
	err := checkAddressActiveByIds(ctx, addrIds)
	if err != nil {
		return err
	}
	return nil
}

func ValidatModifyAddressList(ctx context.Context, req *pb.ModifyAddressListRequest) error {
	addrIds := stringutil.Unique(req.GetAddressId())
	//check addr status is active,only active addr could be add to address list.
	err := checkAddressActiveByIds(ctx, addrIds)
	if err != nil {
		return err
	}
	return nil
}

func ValidateDeleteAddressesList(ctx context.Context, req *pb.DeleteAddressListRequest) error {
	addrIds := stringutil.Unique(req.GetAddressListId())
	//check addr status is active,only active addr could be deleted.
	err := checkAddressListActiveByIds(ctx, addrIds)
	if err != nil {
		return err
	}
	return nil
}
func checkAddressListActiveByIds(ctx context.Context, addrListIds []string) error {
	//check addr is active by ids
	var ActiveAddrLists []*models.AddressList
	ActiveAddrLists, err := rs.GetActiveAddressesListsByIds(ctx, addrListIds)
	if err != nil {
		logger.Errorf(ctx, "Failed to get active address, %+v.", err)
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	if len(ActiveAddrLists) != len(addrListIds) {
		logger.Debugf(ctx, "some address is illegal.")
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalItemInList, addrListIds)
	}
	return nil
}

//****************************************
//Validate Address
//****************************************
func ValidateCreateAddress(ctx context.Context, req *pb.CreateAddressRequest) error {
	//step1:check address existed or not, only address status active does not exist could create.
	//if exist show exist error
	var addrStrs []string
	addrStrs = append(addrStrs, req.GetAddress().GetValue())
	var statuses []string
	statuses = append(statuses, constants.StatusActive)
	var reqDescribeAddresses = &pb.DescribeAddressesRequest{
		Address: addrStrs,
		Status:  statuses,
	}
	_, addrCnt, err := rs.DescribeAddressesWithAddrListId(ctx, reqDescribeAddresses)
	if err != nil {
		logger.Errorf(ctx, "Failed to describe addresses, %+v.", err)
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	if addrCnt != 0 {
		return gerr.New(ctx, gerr.AlreadyExists, gerr.ErrorAlreadyExistResource, addrStrs)
	}

	//step2:verify params, if notify type is email,verify the email fmt
	address := req.GetAddress().GetValue()
	notifyType := req.GetNotifyType().GetValue()

	if notifyType == constants.NotifyTypeEmail {
		err := VerifyEmailFmt(ctx, address)
		if err != nil {
			logger.Errorf(ctx, "Failed to validate address [%s]: %+v", address, err)
			return err
		}
		return nil
	} else if notifyType == constants.NotifyTypeWebsocket {
		return nil
	} else {
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorNotifyType, notifyType)
	}
}

func ValidateDeleteAddresses(ctx context.Context, req *pb.DeleteAddressesRequest) error {
	addrIds := stringutil.Unique(req.GetAddressId())
	//check addr status is active,only active addr could be deleted.
	err := checkAddressActiveByIds(ctx, addrIds)
	if err != nil {
		return err
	}
	return nil
}
func checkAddressActiveByIds(ctx context.Context, addrIds []string) error {
	//check addr is active by ids
	var ActiveAddrs []*models.Address
	ActiveAddrs, err := rs.GetActiveAddressesByIds(ctx, addrIds)
	if err != nil {
		logger.Errorf(ctx, "Failed to get active address, %+v.", err)
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	if len(ActiveAddrs) != len(addrIds) {
		logger.Debugf(ctx, "some address is illegal.")
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalItemInList, addrIds)
	}
	return nil
}

func ValidateModifyAddress(ctx context.Context, req *pb.ModifyAddressRequest) error {
	//if notify type is email, to check email fmt
	address := req.GetAddressDetail().GetValue()
	notifyType := req.GetNotifyType().GetValue()

	if notifyType == constants.NotifyTypeEmail {
		err := VerifyEmailFmt(ctx, address)
		if err != nil {
			logger.Errorf(ctx, "Failed to validate address[%s]: %+v", address, err)
			return gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorValidateFailed, address)
		}
		return nil
	} else {
		return nil
	}
}

//****************************************

//****************************************
//VerifyEmailFmt
func VerifyEmailFmt(ctx context.Context, emailStr string) error {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	reg := regexp.MustCompile(pattern)
	result := reg.MatchString(emailStr)
	if result {
		return nil
	} else {
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalEmailFormat, emailStr)
	}

}

//VerifyPortFmt
func VerifyPortFmt(ctx context.Context, port int32) error {
	if port < 0 || port > 65535 {
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalPort, strconv.Itoa(int(port)))
	} else {
		return nil
	}

}

//VerifyAvailableTimeStr
func VerifyAvailableTimeStr(ctx context.Context, timeStr string) error {
	timeFmt := "15:04:05"
	_, e := time.Parse(timeFmt, timeStr)
	if e != nil {
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalTimeFormat, timeStr)
	}
	return nil
}
