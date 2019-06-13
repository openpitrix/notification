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

	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/gerr"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/pb"
)

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

func ValidateCreateNotificationParams(ctx context.Context, req *pb.CreateNotificationRequest) error {
	//1.validate avaiblable time
	err := validateAvaiblableTime(ctx, req)
	if err != nil {
		return err
	}

	//2.validate address info
	notification := models.NewNotification(
		req.GetContentType().GetValue(),
		req.GetTitle().GetValue(),
		req.GetContent().GetValue(),
		req.GetShortContent().GetValue(),
		req.GetAddressInfo().GetValue(),
		req.GetOwner().GetValue(),
		req.GetExpiredDays().GetValue(),
		req.GetAvailableStartTime().GetValue(),
		req.GetAvailableEndTime().GetValue(),
		req.GetExtra().GetValue(),
	)

	_, decodeMapErr := models.DecodeAddressInfo(notification.AddressInfo)
	if decodeMapErr == nil {
		//2.1. check addressInfo format is like address_info = {"email": ["xxx@abc.com", "xxx@xxx.com"],"websocket": ["system", "huojiao"]}
		err = validateAddressInfo4AddressMap(ctx, notification)
		if err != nil {
			return err
		}
	} else {
		//2.2 check addressInfo format is like address_info = ["adl-xxxx1", "adl-xxxx2"]
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
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalNotificationAddressInfo, notification.AddressInfo)
	}
	return nil
}

func validateAddressInfo4AddressMap(ctx context.Context, notification *models.Notification) error {
	addrInfoMap, err := models.DecodeAddressInfo(notification.AddressInfo)
	if err == nil {
		if len(*addrInfoMap) == 0 {
			logger.Errorf(ctx, "Failed to validate addressInfo,address info is blank: [%s].", notification.AddressInfo)
			return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalNotificationAddressInfo, notification.AddressInfo)
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
				err := models.CheckExtra(ctx, notification.Extra)
				if err != nil {
					return err
				}
			}
		}
		return nil

	} else {
		return err
	}

}

func ValidateCreateAddressParams(ctx context.Context, req *pb.CreateAddressRequest) error {
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
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorValidateFailed)
	}
}

func ValidateModifyAddressParams(ctx context.Context, req *pb.ModifyAddressRequest) error {
	address := req.GetAddress().GetValue()
	notifyType := req.GetNotifyType().GetValue()

	if notifyType == constants.NotifyTypeEmail {
		err := VerifyEmailFmt(ctx, address)
		if err != nil {
			logger.Errorf(ctx, "Failed to validate address [%s]: %+v", address, err)
			return err
		}
		return nil
	} else {
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorValidateFailed)
	}
}

//Email
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

//Port
func VerifyPortFmt(ctx context.Context, port int32) error {
	if port < 0 || port > 65535 {
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalPort, strconv.Itoa(int(port)))
	} else {
		return nil
	}

}

func VerifyAvailableTimeStr(ctx context.Context, timeStr string) error {
	timeFmt := "15:04:05"
	_, e := time.Parse(timeFmt, timeStr)
	if e != nil {
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalTimeFormat, timeStr)
	}
	return nil
}
