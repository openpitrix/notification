package config

import (
	"fmt"
	"github.com/mcuadros/go-defaults"
	"log"
	"sync"
	"time"
)


// Config that contains all of the configuration variables
// that are set up in the environment.
type Config struct {
	DBLogMode         bool   `default:"False"`
	SessionLifeTime   time.Duration

	App struct {
		AppName string `default:"Notification"`
		HostURL string `default:"http://192.168.0.3/"`
		Port    string    `default:":50051"`
		Env     string `default:"DEV"`
		MaxWorkingTasks int `default:"20"`
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
	println("Test Validate start")
	var errorMsg = "config: required field [%v] was not configured!"

	if c.App.HostURL == "" {
		println("Test c.App.HostURL is blank")
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
	println("Test Validate end")
	return nil
}

// Print configuration values to the log. Some user and password fields
// are omitted for security reasons.
func (c *Config) Print() {
	log.Println("----------------------------------")
	log.Println("   Notication Configuration")
	log.Println("----------------------------------")
	log.Println("   DBLogMode:", c.DBLogMode)
	log.Println("   SessionLifeTime:", c.SessionLifeTime)
	log.Println(" ")
	log.Println("   Application HostURL:", c.App.HostURL)
	log.Println("   Application Port:", c.App.Port)
	log.Println("   Application Environment:", c.App.Env)
	log.Println(" ")
	log.Println("   Database Host:", c.Db.Host)
	log.Println("   Database Port:", c.Db.Port)
	log.Println("   Database User:", c.Db.User)
	log.Println("   Database Password:", c.Db.Password)
	log.Println("   Database Database:", c.Db.DatabaseName)
	log.Println("   Database Disable:", c.Db.Disable)
	log.Println("----------------------------------")
}
