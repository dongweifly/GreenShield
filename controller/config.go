package controller

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	GMysqlDB *gorm.DB
)

func Init(address string) error {
	var err error
	GMysqlDB, err = gorm.Open(mysql.New(mysql.Config{
		DriverName: "mysql",
		DSN:        address,
	}))

	return err
}
