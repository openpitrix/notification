// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"strconv"
	"strings"
	"time"

	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/globalcfg"
	"openpitrix.io/notification/pkg/services/notification/service/notification"
	"openpitrix.io/notification/pkg/services/notification/service/task"
	"openpitrix.io/notification/pkg/util/emailutil"

	"openpitrix.io/logger"
	"openpitrix.io/openpitrix/pkg/etcd"
)

type Controller struct {
	queue          *etcd.Queue
	runningTaskIds chan string
	nfIdLast       string
	nfService      notification.Service
	taskService    task.Service
}

func NewController(nfService notification.Service, tasksc task.Service) Controller {
	return Controller{
		queue:          globalcfg.GetInstance().GetEtcd().NewQueue(constants.EmailQueue),
		runningTaskIds: make(chan string),
		nfIdLast:       "",
		nfService:      nfService,
		taskService:    tasksc,
	}
}

func (c *Controller) Serve() {
	go c.ExtractTasks()

	MaxWorkingTasks := constants.MaxWorkingTasks
	for i := 0; i < MaxWorkingTasks; i++ {
		go c.HandleTask(strconv.Itoa(i))
	}
}

func (c *Controller) ExtractTasks() error {
	for {
		nfTaskIdsStr, err := c.queue.Dequeue()
		//taskId := time.Now().Format("2006-01-02 15:04:05")
		//time.Sleep(1 * time.Second)
		if err != nil {
			logger.Errorf(nil, "Failed to dequeue job from etcd queue: %+v", err)
			time.Sleep(3 * time.Second)
			continue
		}

		logger.Infof(nil, "%+v", "Dequeue from etcd queue success,  "+nfTaskIdsStr)
		c.runningTaskIds <- nfTaskIdsStr
	}
	return nil
}

func (c *Controller) HandleTask(handlerNum string) error {
	c.nfIdLast = ""
	for {
		nfTaskIdsStr := <-c.runningTaskIds
		logger.Debugf(nil, time.Now().Format("2006-01-02 15:04:05")+" handlerNum:"+handlerNum+"  Receive:", nfTaskIdsStr)

		ids := strings.Split(nfTaskIdsStr, ",")
		taskId := ids[1]
		nfId := ids[0]

		//logger.Debugf(nil, "test=======handlerNums%d", handlerNum)
		//logger.Debugf(nil, "test=======taskId=s%", taskId)
		//logger.Debugf(nil, "test=======nfId=s%", nfId)

		taskWithNfInfo, err := c.taskService.GetTaskWithNfContentByID(taskId)
		if err != nil {
			logger.Errorf(nil, "Got TaskwithNfContentbyID failed, [%+v]", err)
			return err
		}

		logger.Debugf(nil, "Got TaskwithNfContentbyID successed, : [%+v]", taskWithNfInfo)

		emailAddr := taskWithNfInfo.EmailAddr
		title := taskWithNfInfo.Title
		content := taskWithNfInfo.Content
		err = emailutil.SendMail(emailAddr, title, content)
		if err != nil {
			logger.Warnf(nil, "Send email failed, [%+v]", err)
			//return err
		} else {
			//if send successfully,need to update notification, job and task status.
			_, err = c.taskService.UpdateJobTaskStatus2FinishedById(*taskWithNfInfo)
			if err != nil {
				logger.Errorf(nil, "Update job and task status  to finished failed, [%+v]", err)
				return err
			}
			logger.Debugf(nil, "Update job and task status to finished: [%+v]", taskWithNfInfo)
		}
		//if the nfId is different from nfIdLast,that means the nf including all the tasks is finished.
		//update notification status to finished
		if c.nfIdLast != nfId && c.nfIdLast != "" {
			_, err = c.nfService.UpdateStatus2FinishedById(c.nfIdLast)
			if err != nil {
				logger.Errorf(nil, "Update notification status to finished failed, [%+v]", err)
				return err
			}
		}
		c.nfIdLast = nfId
	}
	return nil
}
