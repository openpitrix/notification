// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package etcd

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"openpitrix.io/logger"
	pkg "openpitrix.io/notification/pkg"
)

func TestEnqueue(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("LocalDevEnv disabled")
	}
	endpoints := []string{"192.168.0.6:12379"}
	//endpoints := []string{"139.198.121.89:12379"}
	e, err := Connect(endpoints, "notification")
	if err != nil {
		t.Fatal(err)
	}
	notificationQueue := e.NewQueue("notification")

	for i := 0; i < 10000; i++ {
		id := fmt.Sprintf("notification_%d", i)
		err := notificationQueue.Enqueue(id)
		if err != nil {
			logger.Errorf(nil, "Failed to dequeue notification from etcd queue: %+v", err)
		}
		logger.Infof(nil, "Enqueue notification [%s] from etcd queue succeed", id)
	}
}

func TestDequeue(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("LocalDevEnv disabled")
	}
	endpoints := []string{"192.168.0.6:12379"}
	e, err := Connect(endpoints, "notification")
	if err != nil {
		t.Fatal(err)
	}

	notificationQueue := e.NewQueue("notification")
	for i := 0; i < 1000; i++ {
		n, err := notificationQueue.Dequeue()
		if err != nil {
			t.Fatal(err)
		}
		logger.Infof(nil, "Got message [%s] from queue, worker number [%d]", n, i)
	}

}

func TestEtcdQueue4OneQueue(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("LocalDevEnv disabled")
	}

	endpoints := []string{"192.168.0.6:12379"}
	e, err := Connect(endpoints, "notification")
	if err != nil {
		t.Fatal(err)
	}

	queue := e.NewQueue(fmt.Sprintf("/test-queue-000"))
	for i := 0; i < 200; i++ {
		err := queue.Enqueue(fmt.Sprintf("msg_%d", i))
		if err != nil {
			t.Fatal(err)
		}
	}
	for i := 0; i < 200; i++ {
		n, err := queue.Dequeue()
		if err != nil {
			t.Fatal(err)
		}
		logger.Infof(nil, "Got message [message=%s] from queue, worker number [%d]", n, i)
	}

	logger.Infof(nil, "Only one Queue, Dequeue speed is about 7 records/s.")
	for {
		time.Sleep(time.Second * 3600)
	}
}

const queueNum = 10

func TestEtcdQueue4MultiQueues(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("LocalDevEnv disabled")
	}

	endpoints := []string{"192.168.0.6:12379"}
	e, err := Connect(endpoints, "notification")
	if err != nil {
		t.Fatal(err)
	}

	testQueues := make([]*Queue, 0)
	for i := 0; i < queueNum; i++ {
		testQueues = append(testQueues, e.NewQueue(fmt.Sprintf("/%s-%d", "testqueue", i)))
	}

	//for i := 0; i < 100; i++ {
	//	go Enqueue2MultiQueue(testQueues )
	//}

	for i := 0; i < 5; i++ {
		go Dequeue2MultiQueue(testQueues)
	}

	logger.Infof(nil, "MultiQueues, Dequeue speed is much more faster.")
	for {
		time.Sleep(time.Second * 3600)
	}
}

func Enqueue2MultiQueue(queues []*Queue) {
	for i := 0; i < 100; i++ {
		queues[rand.Intn(queueNum)].Enqueue(fmt.Sprintf("msg_%d", i))
		logger.Infof(nil, "Enqueue msg [%s] into etcd queue succeed", fmt.Sprintf("msg_%d", i))
	}
}

func Dequeue2MultiQueue(queues []*Queue) {
	for {
		for i := 0; i < queueNum; i++ {
			msg, err := queues[i].Dequeue()
			if err != nil {
				logger.Errorf(nil, "Failed to dequeue notification from etcd queue: %+v", err)
			}
			logger.Infof(nil, "Got message [message=%s] from queue, worker number [%d]", msg, i)
		}
	}
}
