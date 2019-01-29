// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package config

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/koding/multiconfig"

	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/constants"
)

type Config struct {
	Log   LogConfig
	Grpc  GrpcConfig
	Mysql MysqlConfig
	Etcd  EtcdConfig
	Email EmailCfg
}

var instance *Config
var once sync.Once

func GetInstance() *Config {
	once.Do(func() {
		instance = &Config{}
	})
	return instance
}

/*===================================================================================================*/

type LogConfig struct {
	Level string `default:"debug"` // debug, info, warn, error, fatal
}

type GrpcConfig struct {
	ShowErrorCause bool `default:"false"` // show grpc error cause to frontend
}

type EtcdConfig struct {
	Endpoints string `default:"notification-etcd:2379"` // Example: "localhost:2379,localhost:22379,localhost:32379"
	//Endpoints string `default:"192.168.0.7:2379"`
}

type MysqlConfig struct {
	Host string `default:"notification-db"`
	//Host string `default:"192.168.0.10"`
	//Host     string `default:"192.168.0.3"`
	//Port string `default:"13306"`
	Port     string `default:"3306"`
	User     string `default:"root"`
	Password string `default:"password"`
	Database string `default:"notification"`
	Disable  bool   `default:"false"`
	LogMode  bool   `default:"false"`
}

func (m *MysqlConfig) GetUrl() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", m.User, m.Password, m.Host, m.Port, m.Database)
}

type EmailCfg struct {
	Host     string `default:"mail.app-center.cn"`
	Port     int    `default:"25"`
	Username string `default:"openpitrix@app-center.cn"`
	Password string `default:"openpitrix"`
}

/*===================================================================================================*/
func (c *Config) PrintUsage() {
	fmt.Fprintf(os.Stdout, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprint(os.Stdout, "\nSupported environment variables:\n")
	e := newLoader(constants.ServiceName)
	e.PrintEnvs(new(Config))
	fmt.Println("")
}

func (c *Config) GetFlagSet() *flag.FlagSet {
	flag.CommandLine.Usage = c.PrintUsage
	return flag.CommandLine
}

func (c *Config) ParseFlag() {
	c.GetFlagSet().Parse(os.Args[1:])
}

func (c *Config) LoadConf() *Config {
	c.ParseFlag()
	config := instance
	m := &multiconfig.DefaultLoader{}
	m.Loader = multiconfig.MultiLoader(newLoader("notification"))
	m.Validator = multiconfig.MultiValidator(
		&multiconfig.RequiredValidator{},
	)
	err := m.Load(config)
	if err != nil {
		panic(err)
	}

	logger.Infof(nil, "LoadConf: %+v", config)

	return config
}
