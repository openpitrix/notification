// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package plugins

import (
	"context"

	"openpitrix.io/logger"

	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/gerr"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/util/emailutil"
	"openpitrix.io/notification/pkg/util/stringutil"
)

type EmailNotifier struct {
}

func (n *EmailNotifier) Send(ctx context.Context, task *models.Task) error {
	directive, err := models.DecodeTaskDirective(task.Directive)
	if err != nil {
		return err
	}

	//for email msg , use the content with html Format
	contentStruct, err := models.DecodeContent(directive.Content)
	if err != nil {
		logger.Errorf(ctx, "Failed to send notification, content format is not correct, %+v.", err)
		return err
	}
	contentFmtHtml, ok := (*contentStruct)[constants.ContentFmtHtml]
	fmtType := "html"
	if !ok {
		contentFmtHtml = directive.Content
		fmtType = "normal"
	}
	directive.Content = contentFmtHtml

	if directive.AvailableStartTime == "" && directive.AvailableEndTime == "" {
		return emailutil.SendMail(ctx, directive.Address, directive.Title, directive.Content, fmtType)
	} else {
		isOK := stringutil.CheckTimeAvailable(directive.AvailableStartTime, directive.AvailableEndTime)
		if isOK != true {
			logger.Errorf(ctx, "Failed to send notification, time is not available, %+v.", err)
			return gerr.New(nil, gerr.Internal, gerr.ErrorNotAvailableTimeRange, directive.AvailableStartTime, directive.AvailableEndTime)
		}
		return emailutil.SendMail(ctx, directive.Address, directive.Title, directive.Content, fmtType)
	}

}
