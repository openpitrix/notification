// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package etcdutil

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/coreos/etcd/clientv3/namespace"
	"github.com/coreos/etcd/contrib/recipes"
	"openpitrix.io/logger"
	"time"
)

type Etcd struct {
	*clientv3.Client
}

//Connect
func Connect(endpoints []string, prefix string) (*Etcd, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logger.Warnf(nil, "%+v", err)
		return nil, err
	}
	cli.KV = namespace.NewKV(cli.KV, prefix)
	cli.Watcher = namespace.NewWatcher(cli.Watcher, prefix)
	cli.Lease = namespace.NewLease(cli.Lease, prefix)
	return &Etcd{cli}, err
}

func (etcd *Etcd) NewQueue(topic string) *Queue {
	return &Queue{recipe.NewQueue(etcd.Client, topic)}
}

func (etcd *Etcd) NewMutex(key string) (*Mutex, error) {
	session, err := concurrency.NewSession(etcd.Client)
	if err != nil {
		logger.Warnf(nil, "%+v", err)
		return nil, err
	}
	return &Mutex{concurrency.NewMutex(session, key)}, nil
}
