package main

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

type gormLogger struct{}

// Connect to db
func Connect() *gorm.DB {
	var err error
	log.Info("Connecting to db")
	db, err := gorm.Open("mysql", "dancestudio:dancestudio@tcp(127.0.0.1)/dancestudio?charset=utf8&parseTime=True&loc=UTC")
	if err != nil {
		log.Error("Unable to open DB: ", err)
		os.Exit(1)
	}

	db.SetLogger(&gormLogger{})

	return db
}

// Disconnect from db
func Disconnect(db *gorm.DB) {
	log.Info("Disconnecting from db")
	err := db.Close()

	if err != nil {
		log.Error("Unable to open DB: ", err)
	}
}

func (*gormLogger) Print(v ...interface{}) {
	if v[0] == "sql" {
		log.WithFields(log.Fields{"module": "gorm", "type": "sql"}).Print(v[3])
	}
	if v[0] == "log" {
		log.WithFields(log.Fields{"module": "gorm", "type": "log"}).Print(v[2])
	}
}
