package main

import (
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

type gormLogger struct{}

// Config mysql
type Config struct {
	mysqlUser     string
	mysqlPassword string
	mysqlAddress  string
	mysqlDatabase string
	mysqlPort     int
}

var config = Config{
	"dancestudio",
	"dancestudio",
	"127.0.0.1",
	"dancestudio",
	3306,
}

func init() {
	if len(os.Getenv("DANCESTUDIO_MYSQLUSER")) > 0 {
		config.mysqlUser = os.Getenv("DANCESTUDIO_MYSQLUSER")
	}
	if len(os.Getenv("DANCESTUDIO_MYSQLPASSWORD")) > 0 {
		config.mysqlPassword = os.Getenv("DANCESTUDIO_MYSQLPASSWORD")
	}
	if len(os.Getenv("DANCESTUDIO_MYSQLADDRESS")) > 0 {
		config.mysqlAddress = os.Getenv("DANCESTUDIO_MYSQLADDRESS")
	}
	if len(os.Getenv("DANCESTUDIO_MYSQLDB")) > 0 {
		config.mysqlDatabase = os.Getenv("DANCESTUDIO_MYSQLDB")
	}
	if len(os.Getenv("DANCESTUDIO_MYSQLPORT")) > 0 {
		config.mysqlPort, _ = strconv.Atoi(os.Getenv("DANCESTUDIO_MYSQLPORT"))
	}
}

// Connect to db
func Connect() *gorm.DB {
	var err error
	log.Info("Connecting to db")
	db, err := gorm.Open("mysql", config.mysqlUser+":"+config.mysqlPassword+"@tcp("+config.mysqlAddress+":"+strconv.Itoa(config.mysqlPort)+")/"+config.mysqlDatabase+"?charset=utf8&parseTime=True&loc=UTC")
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
