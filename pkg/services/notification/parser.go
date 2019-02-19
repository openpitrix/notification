// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/util/jsonutil"
)

func SplitNotificationIntoTasks(notification *models.Notification) ([]*models.Task, error) {
	addressInfo, err := models.DecodeAddressInfo(notification.AddressInfo)
	if err != nil {
		return nil, err
	}

	var tasks []*models.Task
	for notifyType, addresses := range *addressInfo {
		for _, address := range addresses {
			directive := &models.TaskDirective{
				Address:      address,
				NotifyType:   notifyType,
				ContentType:  notification.ContentType,
				Title:        notification.Title,
				Content:      notification.Content,
				ShortContent: notification.ShortContent,
				ExpiredDays:  notification.ExpiredDays,
			}
			task := models.NewTask(
				notification.NotificationId,
				jsonutil.ToString(directive),
			)

			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}
