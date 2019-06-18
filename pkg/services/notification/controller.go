// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"context"
	"strconv"
	"time"

	lib "openpitrix.io/libqueue"
	q "openpitrix.io/libqueue/queue"
	"openpitrix.io/logger"

	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/plugins"
	rs "openpitrix.io/notification/pkg/services/notification/resource_control"
	"openpitrix.io/notification/pkg/util/ctxutil"
)

type Controller struct {
	runningTaskIds         chan string
	runningNotificationIds chan string
	taskQueue              lib.Topic
	notificationQueue      lib.Topic
}

func NewController() (*Controller, error) {
	cfg := config.GetInstance().LoadConf()
	queueConnStr := cfg.Queue.Addr
	queueType := cfg.Queue.Type

	queueConfigMap := map[string]interface{}{
		"connStr": queueConnStr,
	}

	c, err := q.New(queueType, queueConfigMap)
	if err != nil {
		logger.Errorf(nil, "Failed to connect redis queue: %+v.", err)
		return nil, err
	}

	notificationQueue, _ := c.SetTopic(constants.NotificationTopicPrefix)

	taskQueue, _ := c.SetTopic(constants.NotificationTaskTopicPrefix)
	return &Controller{
		taskQueue:              taskQueue,
		notificationQueue:      notificationQueue,
		runningTaskIds:         make(chan string),
		runningNotificationIds: make(chan string),
	}, nil
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
			logger.Errorf(nil, "Failed to dequeue task from queue: %+v", err)
			time.Sleep(3 * time.Second)
			continue
		}

		logger.Infof(nil, "Dequeue task [%s] from queue succeed", taskId)
		c.runningTaskIds <- taskId
	}
}

func (c *Controller) ExtractNotifications() error {
	for {
		notificationId, err := c.notificationQueue.Dequeue()
		if err != nil {
			logger.Errorf(nil, "Failed to dequeue notification from queue: %+v", err)
			time.Sleep(3 * time.Second)
			continue
		}

		logger.Infof(nil, "Dequeue notification [%s] from queue succeed", notificationId)
		c.runningNotificationIds <- notificationId
	}
}

func (c *Controller) HandleNotification(handlerNum string) {
	for {
		notificationId := <-c.runningNotificationIds
		ctx := ctxutil.AddMessageId(context.Background(), notificationId)
		logger.Debugf(ctx, time.Now().Format("2006-01-02 15:04:05")+" handlerNum:"+handlerNum+"  Receive:"+notificationId)

		//step0: update NF status from pending to sending.
		err := rs.UpdateNotificationsStatus(ctx, []string{notificationId}, constants.StatusSending)
		if err != nil {
			logger.Errorf(ctx, "Update notification status to [sending] failed, [%+v]", err)
			continue
		}

		//when sending,one task has 3 times to retry.
		isNotificationFinished := false

		//taskRetryTimes stores the taskid to retry and the times has retried.【taskId,retryTimes】
		taskRetryTimes := make(map[string]int)
		for {
			if isNotificationFinished {
				break
			}
			noSuccessfulTasks := rs.GetTasksByStatus(ctx,
				notificationId,
				[]string{
					constants.StatusFailed,
					constants.StatusPending,
					constants.StatusSending,
				},
			)

			//step1:check all the taskid for this one notifitication exits not successful status,
			// if not exits, update this one nf status to successful.
			if len(noSuccessfulTasks) == 0 {
				err := rs.UpdateNotificationsStatus(ctx, []string{notificationId}, constants.StatusSuccessful)
				if err != nil {
					logger.Errorf(ctx, "Update notification status to [successful] failed, [%+v]", err)
					continue
				}
				isNotificationFinished = true
			} else {
				//step2: go through all the tasks with not_successful status, if status is failded, retry this task.
				for _, noSuccessfulTask := range noSuccessfulTasks {
					if noSuccessfulTask.Status != constants.StatusFailed {
						continue
					}

					retryTimes, isExist := taskRetryTimes[noSuccessfulTask.TaskId]
					if !isExist {
						retryTimes = 0
					}
					taskRetryTimes[noSuccessfulTask.TaskId] = retryTimes + 1

					//2.1 if the retryTimes for this one task is more than the setting times,
					// update this one notification status to faild.
					if taskRetryTimes[noSuccessfulTask.TaskId] > constants.MaxTaskRetryTimes {
						err := rs.UpdateNotificationsStatus(ctx, []string{notificationId}, constants.StatusFailed)
						if err != nil {
							logger.Errorf(ctx, "Update notification status to [failed] failed, [%+v]", err)
							continue
						}

						//2.1.1 end retry.
						isNotificationFinished = true
					}

					//2.2 retry the task, put this one task back to task queue.
					err := c.taskQueue.Enqueue(noSuccessfulTask.TaskId)
					if err != nil {
						logger.Errorf(nil, "Failed to push task [%s] into queue, error: [%+v]", noSuccessfulTask.TaskId, err)
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

		//step0: update this one task status from pending to sending.
		taskIds := []string{taskId}
		err := rs.UpdateTasksStatus(ctx, taskIds, constants.StatusSending)
		if err != nil {
			logger.Errorf(ctx, "Update task status to [sending] failed, [%+v]", err)
			continue
		}

		//step1: get this task details form db, if not exits in db,just go to handle next task.
		tasks, err := rs.GetTasksByTaskIds(ctx, taskIds)
		if err != nil {
			logger.Errorf(ctx, "Get task failed, [%+v]", err)
			continue
		}
		if len(tasks) == 0 {
			logger.Errorf(ctx, "Get task failed, [%+v]", err)
			continue
		}

		//stpe2: after get task from db just send the task by plugin.
		task := tasks[0]
		notifier, err := plugins.GetNotifier(task)
		if err != nil {
			logger.Errorf(ctx, "Get notifier failed, [%+v]", err)
			continue
		}

		err = notifier.Send(ctx, task)
		if err != nil {
			logger.Errorf(ctx, "Notifier Sends failed, [%+v]", err)

			//if send failed, just update task status to failed in db.
			err := rs.UpdateTasksStatus(ctx, taskIds, constants.StatusFailed)
			if err != nil {
				logger.Errorf(ctx, "Update task status to [failed] failed, [%+v]", err)
				continue
			}
		} else {
			//if send failed, just update task status to successful in db.
			err := rs.UpdateTasksStatus(ctx, taskIds, constants.StatusSuccessful)
			if err != nil {
				logger.Errorf(ctx, "Update task status to [successful] failed, [%+v]", err)
				continue
			}
		}

	}

}
