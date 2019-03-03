// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"context"
	"strconv"
	"time"

	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/etcd"
	"openpitrix.io/notification/pkg/global"
	"openpitrix.io/notification/pkg/plugins"
	rs "openpitrix.io/notification/pkg/services/notification/resource_control"
	"openpitrix.io/notification/pkg/util/ctxutil"
)

type Controller struct {
	taskQueue              *etcd.Queue
	notificationQueue      *etcd.Queue
	runningTaskIds         chan string
	runningNotificationIds chan string
}

func NewController() *Controller {
	return &Controller{
		taskQueue:              global.GetInstance().GetEtcd().NewQueue(constants.NotificationTaskTopicPrefix),
		notificationQueue:      global.GetInstance().GetEtcd().NewQueue(constants.NotificationTopicPrefix),
		runningTaskIds:         make(chan string),
		runningNotificationIds: make(chan string),
	}
}

func (c *Controller) Serve() {
	go c.ExtractTasks()
	go c.ExtractNotifications()

	for i := 0; i < constants.MaxWorkingTasks; i++ {
		go c.HandleTask(strconv.Itoa(i))
	}
	for i := 0; i < constants.MaxWorkingNotifications; i++ {
		go c.HandleNotification(strconv.Itoa(i))
	}
}

func (c *Controller) ExtractTasks() error {
	for {
		taskId, err := c.taskQueue.Dequeue()
		if err != nil {
			logger.Errorf(nil, "Failed to dequeue task from etcd queue: %+v", err)
			time.Sleep(3 * time.Second)
			continue
		}

		logger.Debugf(nil, "Dequeue task [%s] from etcd queue succeed", taskId)
		c.runningTaskIds <- taskId
	}
}

func (c *Controller) ExtractNotifications() error {
	for {
		notificationId, err := c.notificationQueue.Dequeue()
		if err != nil {
			logger.Errorf(nil, "Failed to dequeue notification from etcd queue: %+v", err)
			time.Sleep(3 * time.Second)
			continue
		}

		logger.Debugf(nil, "Dequeue notification [%s] from etcd queue succeed", notificationId)
		c.runningNotificationIds <- notificationId
	}
}

func (c *Controller) HandleNotification(handlerNum string) {
	for {
		notificationId := <-c.runningNotificationIds
		ctx := ctxutil.AddMessageId(context.Background(), notificationId)

		logger.Debugf(ctx, time.Now().Format("2006-01-02 15:04:05")+" handlerNum:"+handlerNum+"  Receive:"+notificationId)

		err := rs.UpdateNotificationStatus(notificationId, constants.StatusSending)
		if err != nil {
			logger.Errorf(ctx, "Update notification status to [sending] failed, [%+v]", err)
			continue
		}

		isNotificationFinished := false
		taskRetryTimes := make(map[string]int)
		for {
			if isNotificationFinished {
				break
			}
			noSuccessfulTasks := rs.GetTasksByStatus(
				notificationId,
				[]string{
					constants.StatusFailed,
					constants.StatusPending,
					constants.StatusSending,
				},
			)
			if len(noSuccessfulTasks) == 0 {
				err := rs.UpdateNotificationStatus(notificationId, constants.StatusSuccessful)
				if err != nil {
					logger.Errorf(ctx, "Update notification status to [successful] failed, [%+v]", err)
					continue
				}
				isNotificationFinished = true
			} else {
				for _, noSuccessfulTask := range noSuccessfulTasks {
					if noSuccessfulTask.Status != constants.StatusFailed {
						continue
					}
					retryTimes, isExist := taskRetryTimes[noSuccessfulTask.TaskId]
					if !isExist {
						retryTimes = 0
					}
					taskRetryTimes[noSuccessfulTask.TaskId] = retryTimes + 1
					if taskRetryTimes[noSuccessfulTask.TaskId] > constants.MaxTaskRetryTimes {
						err := rs.UpdateNotificationStatus(notificationId, constants.StatusFailed)
						if err != nil {
							logger.Errorf(ctx, "Update notification status to [failed] failed, [%+v]", err)
							continue
						}
						isNotificationFinished = true
					}

					err := c.taskQueue.Enqueue(noSuccessfulTask.TaskId)
					if err != nil {
						logger.Errorf(nil, "Failed to push task [%s] into etcd, error: [%+v]", noSuccessfulTask.TaskId, err)
						continue
					}
				}
			}
			time.Sleep(3 * time.Second)
		}
	}
}

func (c *Controller) HandleTask(handlerNum string) {
	for {
		taskId := <-c.runningTaskIds
		ctx := ctxutil.AddMessageId(context.Background(), taskId)

		logger.Debugf(ctx, time.Now().Format("2006-01-02 15:04:05")+" handlerNum:"+handlerNum+"  Receive:"+taskId)

		err := rs.UpdateTaskStatus(taskId, constants.StatusSending)
		if err != nil {
			logger.Errorf(ctx, "Update task status to [sending] failed, [%+v]", err)
			continue
		}

		task, err := rs.GetTask(taskId)
		if err != nil {
			logger.Errorf(ctx, "Get task failed, [%+v]", err)
			continue
		}

		notifier, err := plugins.GetNotifier(task)
		if err != nil {
			logger.Errorf(ctx, "Get notifier failed, [%+v]", err)
			continue
		}

		err = notifier.Send(ctx, task)
		if err != nil {
			logger.Errorf(ctx, "Notifier Send failed, [%+v]", err)
			err = rs.UpdateTaskStatus(taskId, constants.StatusFailed)
			if err != nil {
				logger.Errorf(ctx, "Update task status to [failed] failed, [%+v]", err)
				continue
			}
		} else {
			err = rs.UpdateTaskStatus(taskId, constants.StatusSuccessful)
			if err != nil {
				logger.Errorf(ctx, "Update task status to [successful] failed, [%+v]", err)
				continue
			}
		}
	}
}
