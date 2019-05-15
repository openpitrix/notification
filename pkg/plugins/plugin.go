// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package plugins

import (
	"context"
	"fmt"

	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/models"
)

type NotifyInterface interface {
	Send(ctx context.Context, task *models.Task) error
}

func GetNotifier(task *models.Task) (NotifyInterface, error) {
	taskDirective, err := models.DecodeTaskDirective(task.Directive)
	if err != nil {
		return nil, err
	}
	switch taskDirective.NotifyType {
	case constants.NotifyTypeEmail:
		return new(EmailNotifier), nil
	default:
		return nil, fmt.Errorf("unsupported notify type [%s]", taskDirective.NotifyType)
	}
}
