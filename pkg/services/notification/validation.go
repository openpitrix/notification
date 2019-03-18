// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"context"
	"regexp"
	"strconv"

	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/gerr"
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
