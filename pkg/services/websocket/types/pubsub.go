package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/namespace"
	"github.com/go-redis/redis"
	"openpitrix.io/logger"
	"strings"
	"time"
)

type PubsubClient interface {
}

type EtcdClient struct {
	*clientv3.Client
}

type RedisClient struct {
	*redis.Client
}

func New(pubsubType string, configMap map[string]interface{}) (PubsubClient, error) {
	if configMap == nil {
		return nil, fmt.Errorf("not provide queue configuration info.")
	}

	switch pubsubType {
	case "etcd":
		cfg := LoadConf4Etcd(configMap)
		if cfg.ConnStr == "" {
			return nil, errors.New("not provide ConnStr parameter.")
		}

		var dialTimeout time.Duration = (time.Duration(5) * 1000) * time.Millisecond
		if cfg.DialTimeoutSecond != 0 {
			dialTimeout = (time.Duration(cfg.DialTimeoutSecond) * 1000) * time.Millisecond
		}
		cli, err := clientv3.New(clientv3.Config{
			Endpoints:   strings.Split(cfg.ConnStr, ","),
			DialTimeout: dialTimeout,
		})
		if err != nil {
			return nil, err
		}
		cli.KV = namespace.NewKV(cli.KV, "")
		cli.Watcher = namespace.NewWatcher(cli.Watcher, "")
		cli.Lease = namespace.NewLease(cli.Lease, "")

		return EtcdClient{cli}, err

	case "redis":
		cfg := LoadConf4Redis(configMap)

		if cfg.ConnStr == "" {
			return nil, errors.New("not provide ConnStr parameter.")
		}

		options, err := redis.ParseURL(cfg.ConnStr)
		if err != nil {
			return nil, err
		}

		if cfg.PoolSize != 0 {
			options.PoolSize = cfg.PoolSize
		}

		if cfg.MinIdleConns != 0 {
			options.MinIdleConns = cfg.MinIdleConns
		}

		cli := redis.NewClient(options)
		return RedisClient{cli}, nil
	default:
		return nil, fmt.Errorf("unsupported queueType [%s]", pubsubType)
	}
}

type RedisConfig struct {
	ConnStr      string
	PoolSize     int
	MinIdleConns int
}

type EtcdConfig struct {
	ConnStr           string
	DialTimeoutSecond int
}

func LoadConf4Redis(configMap map[string]interface{}) *RedisConfig {
	mjson, _ := json.Marshal(configMap)
	mString := string(mjson)
	logger.Debugf(nil, "mString:%s", mString)

	var config RedisConfig
	data := []byte(mString)
	err := json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println(err)
	}
	return &config
}

func LoadConf4Etcd(configMap map[string]interface{}) *EtcdConfig {
	mjson, _ := json.Marshal(configMap)
	mString := string(mjson)
	logger.Debugf(nil, "mString:%s", mString)

	var config EtcdConfig
	data := []byte(mString)
	err := json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println(err)
	}
	return &config
}
