package config

import (
	"flag"
	"fmt"
	"github.com/koding/multiconfig"
	"github.com/mcuadros/go-defaults"
	"openpitrix.io/logger"
	"os"
	"sync"
)

// Config that contains all of the configuration variables
// that are set up in the environment.
type Config struct {
	//AppLogMode string `default:"info"`
	AppLogMode string `default:"debug"`
	App struct {
		AppName         string `default:"Notification"`
		Host         string `default:"192.168.0.3"`
		//Host         string `default:"localhost"`
		Port            string `default:":50051"`
		Env             string `default:"DEV"`
		MaxWorkingTasks int    `default:"5"` //default:"20
	}

	Db struct {
		Host         string `default:"192.168.0.10"`
		Port         string `default:"13306"`
		User         string `default:"root"`
		Password     string `default:"password"`
		DatabaseName string `default:"notification"`
		Disable      bool   `default:"true"`
		DBLogMode    bool   `default:"true"`
	}

	Etcd struct {
		//Endpoints string `default:"192.168.0.7:2379,192.168.0.8:2379,192.168.0.6:2379"` // Example: "localhost:2379,localhost:22379,localhost:32379"  or default:"openpitrix-etcd:2379
		Endpoints  string `default:"192.168.0.7:2379"`
		Etcdprefix string `default:"nf_"`
		Etcdtopic  string `default:"task"`
	}

	Email struct {
	 	EmailHost  string `default:"mail.app-center.cn"`
		EmailPort int `default:"25"`
		EmailUsername  string `default:"openpitrix@app-center.cn"`
		EmailPassword string `default:"openpitrix"`
	}
}

var instance *Config
var once sync.Once

func GetInstance() *Config {
	once.Do(func() {
		instance = &Config{}
	})
	return instance
}

func (c *Config) InitCfg() {
	defaults.SetDefaults(instance)
}


// Validate checks if the most important fields are set and are not empty
// values.
func (c *Config) Validate() error {
	var errorMsg = "config: required field [%v] was not configured!"

	logger.Infof(nil, "%+v", "==============Validate Start===================================")
	if c.AppLogMode == "" {
		return fmt.Errorf(errorMsg, "App.AppLogMode")
	}

	logger.Infof(nil, "%+v", "-------------App cfg---------------------")
	if c.App.Port == "" {
		return fmt.Errorf(errorMsg, "App.Port")
	}
	if c.App.AppName == "" {
		return fmt.Errorf(errorMsg, "App.AppName")
	}
	if c.App.Env == "" {
		return fmt.Errorf(errorMsg, "App.Env")
	}
	if c.App.Host == "" {
		return fmt.Errorf(errorMsg, "App.Host")
	}
	if c.App.MaxWorkingTasks == 0 {
		return fmt.Errorf(errorMsg, "App.MaxWorkingTasks")
	}

	logger.Infof(nil, "%+v", "-------------Db cfg---------------------")
	if c.Db.Port == "" {
		return fmt.Errorf(errorMsg, "Db.Port")
	}
	if c.Db.Host == "" {
		return fmt.Errorf(errorMsg, "Db.Host")
	}
	if c.Db.DatabaseName == "" {
		return fmt.Errorf(errorMsg, "Db.DatabaseName")
	}
	if c.Db.User == "" {
		return fmt.Errorf(errorMsg, "Db.User")
	}
	if c.Db.Password == "" {
		return fmt.Errorf(errorMsg, "Db.Password")
	}
	if c.Db.Disable == false {
		return fmt.Errorf(errorMsg, "Db.Disable")
	}
	//if c.Db.DBLogMode == false{
	//	return fmt.Errorf(errorMsg, "Db.DBLogMode")
	//}

	logger.Infof(nil, "%+v", "-------------Etcd cfg---------------------")
	if c.Etcd.Endpoints == "" {
		return fmt.Errorf(errorMsg, "Etcd.Endpoints")
	}
	if c.Etcd.Etcdprefix == "" {
		return fmt.Errorf(errorMsg, "Etcd.Etcdprefix")
	}
	if c.Etcd.Etcdtopic == "" {
		return fmt.Errorf(errorMsg, "Etcd.Etcdtopic")
	}
	logger.Infof(nil, "%+v", "==============Validate Start===================================")

	return nil
}

// Print configuration values to the log. Some user and password fields
// are omitted for security reasons.
func (c *Config) Print() {
	logger.Infof(nil, "%+v", "===============Print Start==================================")
	logger.Infof(nil, "%+v", "Notication Configuration")
	logger.Infof(nil, "%+v", "-------------cfg---------------------")
	logger.Infof(nil, "c.AppLogMode:%+v", c.AppLogMode)
	logger.Infof(nil, "%+v", "-------------App cfg---------------------")
	logger.Infof(nil, "c.App.Port:%+v", c.App.Port)
	logger.Infof(nil, "c.App.MaxWorkingTasks:%+v", c.App.MaxWorkingTasks)
	logger.Infof(nil, "c.App.Host:%+v", c.App.Host)
	logger.Infof(nil, "c.App.Env:%+v", c.App.Env)
	logger.Infof(nil, "c.App.AppName:%+v", c.App.AppName)

	logger.Infof(nil, "%+v", "-------------Db cfg---------------------")
	logger.Infof(nil, "c.Db.Port:%+v", c.Db.Port)
	logger.Infof(nil, "c.Db.DBLogMode:%+v", c.Db.DBLogMode)
	logger.Infof(nil, "c.Db.User:%+v", c.Db.User)
	logger.Infof(nil, "c.Db.Password:%+v", c.Db.Password)
	logger.Infof(nil, "c.Db.DatabaseName:%+v", c.Db.DatabaseName)
	logger.Infof(nil, "c.Db.Host:%+v", c.Db.Host)
	logger.Infof(nil, "c.Db.Disable:%+v", c.Db.Disable)

	logger.Infof(nil, "%+v", "-------------Etcd cfg---------------------")
	logger.Infof(nil, "c.Etcd.Etcdprefix:%+v", c.Etcd.Etcdprefix)
	logger.Infof(nil, "c.Etcd.Etcdtopic:%+v", c.Etcd.Etcdtopic)
	logger.Infof(nil, "c.Etcd.Endpoints:%+v", c.Etcd.Endpoints)

	logger.Infof(nil, "%+v", "===============Print End==================================")
}


func PrintUsage() {
	fmt.Fprintf(os.Stdout, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprint(os.Stdout, "\nSupported environment variables:\n")
	e := newLoader("notification")
	e.PrintEnvs(new(Config))
	fmt.Println("")
}

func GetFlagSet() *flag.FlagSet {
	flag.CommandLine.Usage = PrintUsage
	return flag.CommandLine
}

func ParseFlag() {
	GetFlagSet().Parse(os.Args[1:])
}

func  (c *Config)LoadConf()   {
	ParseFlag()
	config:=instance
	m := &multiconfig.DefaultLoader{}
	m.Loader = multiconfig.MultiLoader(newLoader("notification"))
	m.Validator = multiconfig.MultiValidator(
		&multiconfig.RequiredValidator{},
	)
	err := m.Load(config)
	if err != nil {
		logger.Criticalf(nil, "Failed to load config: %+v", err)
		panic(err)
	}
	logger.Debugf(nil, "LoadConf: %+v", config)
}