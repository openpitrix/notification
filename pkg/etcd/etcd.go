// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package etcd

import (
	"strconv"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/namespace"

	"openpitrix.io/notification/pkg/config"
)

type Etcd struct {
	*clientv3.Client
}

func Connect(endpoints []string, prefix string) (*Etcd, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	cli.KV = namespace.NewKV(cli.KV, prefix)
	cli.Watcher = namespace.NewWatcher(cli.Watcher, prefix)
	cli.Lease = namespace.NewLease(cli.Lease, prefix)
	return &Etcd{cli}, err
}

func GetQueueNum() int {
	cfg := config.GetInstance()
	queueNum, err := strconv.ParseInt(cfg.Etcd.QueueNum, 10, 0)
	if err != nil {
		queueNum = 100
	}

	if queueNum < 10 {
		queueNum = 10
	}

	if queueNum > 1000 {
		queueNum = 1000
	}

	return int(queueNum)
}
