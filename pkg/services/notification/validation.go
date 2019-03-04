// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"context"
	"regexp"
	"strconv"

	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/openpitrix/pkg/gerr"
)

func ValidateSetServiceConfigParams(ctx context.Context, req *pb.ServiceConfig) error {
	email := req.GetEmailServiceConfig().GetEmail().GetValue()
	err := VerifyEmailFmt(ctx, email)
	if err != nil {
		logger.Errorf(ctx, "Failed to validateSetServiceConfigParams [%s], %+v", email, err)
		return err
	}

	//displayEmail := req.GetEmailServiceConfig().GetDisplayEmail().GetValue()
	//err = VerifyEmailFmt(ctx, displayEmail)
	//if err != nil {
	//	logger.Errorf(ctx, "Failed to validateSetServiceConfigParams [%s], %+v", displayEmail, err)
	//	return err
	//}

	portStr := req.GetEmailServiceConfig().GetPort().GetValue()
	portNum, err := strconv.ParseInt(portStr, 10, 64)
	err = VerifyPortFmt(ctx, portNum)
	if err != nil {
		logger.Errorf(ctx, "Failed to validateSetServiceConfigParams [%s], %+v", portStr, err)
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
