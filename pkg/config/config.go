package config

import (
	"flag"
	"fmt"
	"github.com/koding/multiconfig"
	"openpitrix.io/logger"
	"os"
	"sync"
)

// Config that contains all of the configuration variables
// that are set up in the environment.
type Config struct {
	App   Appcfg
	Db    Dbcfg
	Etcd  Etcdcfg
	Email Emailcfg
}

type Appcfg struct {
	Name string `default:"Notification"`
	//Host    string `default:"192.168.0.3"`
	Host       string `default:"127.0.0.1"`
	Port       string `default:":50051"`
	Env        string `default:"DEV"`
	Maxtasks   int    `default:"5"`
	Applogmode string `default:"debug"`
}

type Dbcfg struct {
	Host     string `default:"192.168.0.10"`
	Port     string `default:"13306"`
	User     string `default:"root"`
	Password string `default:"password"`
	Dbname   string `default:"notification"`
	Disable  bool   `default:"true"`
	Logmode  bool   `default:"true"`
}

type Etcdcfg struct {
	//Endpoints string `default:"192.168.0.7:2379,192.168.0.8:2379,192.168.0.6:2379"` // Example: "localhost:2379,localhost:22379,localhost:32379"  or default:"openpitrix-etcd:2379
	Endpoints string `default:"192.168.0.7:2379"`
	Prefix    string `default:"nf_"`
	Topic     string `default:"task"`
}

type Emailcfg struct {
	Host     string `default:"mail.app-center.cn"`
	Port     int    `default:"25"`
	Username string `default:"openpitrix@app-center.cn"`
	Password string `default:"openpitrix"`
}

var instance *Config
var once sync.Once

func GetInstance() *Config {
	once.Do(func() {
		instance = &Config{}
	})
	return instance
}

func (c *Config) PrintUsage() {
	flag.PrintDefaults()
	fmt.Fprint(os.Stdout, "\nSupported environment variables:\n")
	e := newLoader("nf")
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

var profilingServerStarted = false

func (c *Config) LoadConf() *Config {
	c.ParseFlag()
	config := instance
	//config := new(Config)
	m := &multiconfig.DefaultLoader{}
	m.Loader = multiconfig.MultiLoader(newLoader("nf"))
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
