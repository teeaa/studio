package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

var db *gorm.DB

func main() {
	log.SetLevel(log.DebugLevel)
	fmt.Println("Dance studio Go server")

	srv := startServer()
	defer srv.Close()

	waitForExit()
}

func waitForExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
