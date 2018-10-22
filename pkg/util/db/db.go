package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
)


// Server defines the web server structure.
type Database struct {
}


// createDatabaseConn creates a new GORM database with the specified database
// configuration.
func (s *Database) createDatabaseConnection() (*gorm.DB, error) {
	var (
		db               *gorm.DB
		err              error
		dbCfg            = s.cfg.DB
		connectionString = fmt.Sprintf(
			"host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
			dbCfg.Host,
			dbCfg.Port,
			dbCfg.User,
			dbCfg.Password,
			dbCfg.Name,
		)
	)

	db, err = gorm.Open("mysql", connectionString)

	if err != nil {
		return nil, err
	}

	err = db.DB().Ping()

	if err != nil {
		return nil, err
	}

	db.DB().SetMaxIdleConns(10)
	return db, nil
}
