// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"context"
	"regexp"
	"strconv"

	"openpitrix.io/notification/pkg/constants"

	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/gerr"
	"openpitrix.io/notification/pkg/pb"
)

func ValidateSetServiceConfigParams(ctx context.Context, req *pb.ServiceConfig) error {
	email := req.GetEmailServiceConfig().GetEmail().GetValue()
	err := VerifyEmailFmt(ctx, email)
	if err != nil {
		logger.Errorf(ctx, "Failed to validate SetServiceConfig Params [%s], %+v", email, err)
		return err
	}

	portStr := req.GetEmailServiceConfig().GetPort().GetValue()
	portNum, err := strconv.ParseInt(portStr, 10, 64)
	err = VerifyPortFmt(ctx, portNum)
	if err != nil {
		logger.Errorf(ctx, "Failed to validate SetServiceConfig Params [%s], %+v", portStr, err)
		return err
	}
	return nil
}

//Email
func VerifyEmailFmt(ctx context.Context, emailStr string) error {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	reg := regexp.MustCompile(pattern)
	result := reg.MatchString(emailStr)
	if result {
		return nil
	} else {
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorValidateFailed, emailStr)
	}

}

//Port
func VerifyPortFmt(ctx context.Context, port int64) error {
	if port < 0 || port > 65535 {
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorValidateFailed, string(port))
	} else {
		return nil
	}

}

func ValidateCreateAddressParams(ctx context.Context, req *pb.CreateAddressRequest) error {
	address := req.GetAddress().GetValue()
	notifyType := req.GetNotifyType().GetValue()

	if notifyType == constants.NotifyTypeEmail {
		err := VerifyEmailFmt(ctx, address)
		if err != nil {
			logger.Errorf(ctx, "Failed to Validate CreateAddress Params [%s], %+v", address, err)
			return err
		}
	}
	return nil
}

func ValidateModifyAddressParams(ctx context.Context, req *pb.ModifyAddressRequest) error {
	address := req.GetAddress().GetValue()
	notifyType := req.GetNotifyType().GetValue()

	if notifyType == constants.NotifyTypeEmail {
		err := VerifyEmailFmt(ctx, address)
		if err != nil {
			logger.Errorf(ctx, "Failed to Validate ModifyAddress Params [%s], %+v", address, err)
			return err
		}
	}
	return nil
}
