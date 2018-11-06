package config

import (
	"fmt"
	"github.com/mcuadros/go-defaults"
	"sync"
	"time"
	"openpitrix.io/logger"
)


// Config that contains all of the configuration variables
// that are set up in the environment.
type Config struct {
	AppLogMode        string `default:"debug"`
	DBLogMode         bool   `default:"False"`
	SessionLifeTime   time.Duration

	App struct {
		AppName string `default:"Notification"`
		HostURL string `default:"http://192.168.0.3/"`
		Port    string    `default:":50051"`
		Env     string `default:"DEV"`
		MaxWorkingTasks int `default:"5"`    //default:"20
	}

	Db struct {
		Host     string `default:"192.168.0.10"`
		Port     string `default:"13306"`
		User     string `default:"root"`
		Password string `default:"password"`
		DatabaseName string `default:"notification"`
		Disable  bool   `default:"true"`
	}

	Etcd struct{
		//Endpoints string `default:"192.168.0.7:2379,192.168.0.8:2379,192.168.0.6:2379"` // Example: "localhost:2379,localhost:22379,localhost:32379"  or default:"openpitrix-etcd:2379
		Endpoints string `default:"192.168.0.7:2379"`
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


func (c *Config) InitCfg()   {
	defaults.SetDefaults(instance)
}


// Validate checks if the most important fields are set and are not empty
// values.
func (c *Config) Validate() error {
	var errorMsg = "config: required field [%v] was not configured!"

	if c.App.HostURL == "" {
		return fmt.Errorf(errorMsg, "App.HostURL")
	}

	if c.Db.Host == "" {
		return fmt.Errorf(errorMsg, "Db.Host")
	}

	if c.Db.Port == "" {
		return fmt.Errorf(errorMsg, "Db.Port")
	}

	if c.Db.User == "" {
		return fmt.Errorf(errorMsg, "Db.User")
	}

	if c.Db.Password == "" {
		return fmt.Errorf(errorMsg, "Db.Password")
	}

	if c.Db.DatabaseName == "" {
		return fmt.Errorf(errorMsg, "Db.Database")
	}

	if c.Db.Disable == false {
		return fmt.Errorf(errorMsg, "Db.Disable")
	}
	return nil
}

// Print configuration values to the log. Some user and password fields
// are omitted for security reasons.
func (c *Config) Print() {
	logger.Infof(nil,"%+v","----------------------------------")
	logger.Infof(nil,"%+v","   Notication Configuration")
	logger.Infof(nil,"%+v","----------------------------------")
	logger.Infof(nil, "DBLogMode: %+v", c.DBLogMode )
	logger.Infof(nil, "SessionLifeTime: %+v", c.SessionLifeTime )
	logger.Infof(nil,"%+v"," ")
	logger.Infof(nil, "Application HostURL: %+v", c.App.HostURL )
	logger.Infof(nil, "Application Port: %+v", c.App.Port )
	logger.Infof(nil, "Application Environment: %+v", c.App.Env )
	logger.Infof(nil, "---" )
	logger.Infof(nil, "Database Host: %+v",  c.Db.Host )
	logger.Infof(nil, "Database User: %+v",  c.Db.User )
	logger.Infof(nil, "Database Password: %+v",   c.Db.Password )
	logger.Infof(nil, "Database Database: %+v",  c.Db.DatabaseName )
	logger.Infof(nil, "Database Disable:%+v",  c.Db.Disable)
	logger.Infof(nil,"----------------------------------")
}
