package main

import (
	"fmt"
	"log"
	"os"
)

func InitLogger() {
	_ = os.Mkdir("./logs", 0777)
	logFile, err := os.OpenFile("./logs/ContainerManager.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Open log file failed, err:", err)
		panic(err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
}
