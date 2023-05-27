package model

import (
	"admin/conf"
	"admin/model/model"
	"admin/model/query"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var Q *query.Query

func DBInit() {
	url := fmt.Sprintf("%s:%s@(%s:%d)/%s",
		conf.Config.MySQL.Username,
		conf.Config.MySQL.Password,
		conf.Config.MySQL.Address,
		conf.Config.MySQL.Port,
		conf.Config.MySQL.DBName,
	)
	conn := mysql.Open(url)
	ormConn, err := gorm.Open(conn)
	if err != nil {
		log.Fatal(fmt.Errorf("[MySQL]: Init Database Connection Failed: %s", err.Error()))
	}
	Q = query.Use(ormConn)
}

func TestDBInit() error {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&model.UserUser{}, &model.Function{}, &model.Trigger{}, &model.Template{})
	if err != nil {
		return err
	}

	Q = query.Use(db)
	return nil
}
