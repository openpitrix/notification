// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package plugins

import (
	"context"

	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/util/emailutil"
)

type EmailNotifier struct {
}

func (n *EmailNotifier) Send(ctx context.Context, task *models.Task) error {
	directive, err := models.DecodeTaskDirective(task.Directive)
	if err != nil {
		return err
	}
	return emailutil.SendMail(ctx, directive.Address, directive.Title, directive.Content)
}
