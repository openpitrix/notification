// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"context"
	"strconv"
	"strings"
	"time"

	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/globalcfg"
	"openpitrix.io/notification/pkg/util/emailutil"

	"openpitrix.io/logger"
	"openpitrix.io/openpitrix/pkg/etcd"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
)

type Controller struct {
	taskQueue      *etcd.Queue
	jobQueue       *etcd.Queue
	runningTaskIds chan string
	runningJobIds  chan string
}

func NewController() *Controller {
	return &Controller{
		taskQueue:      globalcfg.GetInstance().GetEtcd().NewQueue(constants.NotificationTaskTopic),
		jobQueue:       globalcfg.GetInstance().GetEtcd().NewQueue(constants.NotificationJobTopic),
		runningTaskIds: make(chan string),
		runningJobIds:  make(chan string),
	}
}

func (c *Controller) Serve() {
	go c.ExtractTasks()
	go c.ExtractJobs()

	for i := 0; i < constants.MaxWorkingTasks; i++ {
		go c.HandleTask(strconv.Itoa(i))
	}
	for i := 0; i < constants.MaxWorkingJobs; i++ {
		go c.HandleJob(strconv.Itoa(i))
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

func (c *Controller) ExtractJobs() error {
	for {
		jobId, err := c.jobQueue.Dequeue()
		if err != nil {
			logger.Errorf(nil, "Failed to dequeue job from etcd queue: %+v", err)
			time.Sleep(3 * time.Second)
			continue
		}

		logger.Debugf(nil, "Dequeue job [%s] from etcd queue succeed", jobId)
		c.runningJobIds <- jobId
	}
}

func (c *Controller) HandleJob(handlerNum string) error {
	ctx := context.Background()
	for {
		nfIdAndJobId := <-c.runningJobIds
		logger.Debugf(ctx, time.Now().Format("2006-01-02 15:04:05")+" handlerNum:"+handlerNum+"  Receive:"+nfIdAndJobId)

		isJobFinished := false
		taskRetryTimes := make(map[string]int)
		sp := strings.Split(nfIdAndJobId, ",")
		if len(sp) != 2 {
			logger.Criticalf(ctx, "Failed to handle job [%s]", nfIdAndJobId)
			continue
		}
		nfId := sp[0]
		jobId := sp[1]

		ctx = ctxutil.AddMessageId(ctx, nfId)
		ctx = ctxutil.AddMessageId(ctx, jobId)

		for {
			if isJobFinished {
				break
			}
			noSuccessfulTasks := GetStatusTasks(jobId,
				[]string{
					constants.StatusFailed,
					constants.StatusNew,
					constants.StatusSending,
				},
			)
			if len(noSuccessfulTasks) == 0 {
				err := UpdateJobStatus(jobId, constants.StatusSuccessful)
				if err != nil {
					logger.Errorf(ctx, "Update job status to successful failed, [%+v]", err)
				}

				err = UpdateNfStatus(jobId, constants.StatusSuccessful)
				if err != nil {
					logger.Errorf(ctx, "Update nf status to successful failed, [%+v]", err)
				}
				isJobFinished = true
			} else {
				for _, noSuccessfulTask := range noSuccessfulTasks {
					if noSuccessfulTask.Status != constants.StatusFailed {
						continue
					}
					retryTimes, isExist := taskRetryTimes[noSuccessfulTask.TaskID]
					if !isExist {
						retryTimes = 0
					}
					taskRetryTimes[noSuccessfulTask.TaskID] = retryTimes + 1
					if taskRetryTimes[noSuccessfulTask.TaskID] > constants.MaxTaskRetryTimes {
						err := UpdateJobStatus(jobId, constants.StatusFailed)
						if err != nil {
							logger.Errorf(ctx, "Update job status to failed failed, [%+v]", err)
						}

						err = UpdateNfStatus(jobId, constants.StatusFailed)
						if err != nil {
							logger.Errorf(ctx, "Update nf status to failed failed, [%+v]", err)
						}
						isJobFinished = true
					}

					err := c.taskQueue.Enqueue(noSuccessfulTask.TaskID)
					if err != nil {
						logger.Errorf(nil, "Failed to push task [%s] into etcd, error: [%+v]", noSuccessfulTask.TaskID, err)
					}
				}
			}
			time.Sleep(3 * time.Second)
		}
	}
}

func (c *Controller) HandleTask(handlerNum string) error {
	ctx := context.Background()
	for {
		taskId := <-c.runningTaskIds
		logger.Debugf(ctx, time.Now().Format("2006-01-02 15:04:05")+" handlerNum:"+handlerNum+"  Receive:"+taskId)

		ctx = ctxutil.AddMessageId(ctx, taskId)

		taskWithNfInfo, err := GetTaskWithNfInfo(taskId)
		if err != nil {
			logger.Criticalf(ctx, "Get task failed, [%+v]", err)
		}

		logger.Debugf(ctx, "Get task succeed: [%+v]", taskWithNfInfo)

		emailAddr := taskWithNfInfo.EmailAddr
		title := taskWithNfInfo.Title
		content := taskWithNfInfo.Content
		err = emailutil.SendMail(emailAddr, title, content)
		if err != nil {
			logger.Errorf(ctx, "Send email to [%s] failed, [%+v]", emailAddr, err)
			err = UpdateTaskStatus(taskWithNfInfo.TaskID, constants.StatusFailed)
			if err != nil {
				logger.Errorf(ctx, "Update task status to failed failed, [%+v]", err)
			}
		} else {
			// if send successfully, need to update notification, job and task status.
			err = UpdateTaskStatus(taskWithNfInfo.TaskID, constants.StatusSuccessful)
			if err != nil {
				logger.Errorf(ctx, "Update task status to successful failed, [%+v]", err)
			}
		}
	}
}
