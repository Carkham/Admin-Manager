package main

import (
	"admin/conf"
	"admin/model"
	"admin/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	InitLogger()
	model.DBInit()
	utils.InitRedisClient()
	utils.KubeClientInit(conf.Config.K8S.ConfigPath)
	r := gin.Default()
	customRouter(r)
	listenUrl := fmt.Sprintf("0.0.0.0:%d", conf.Config.Service.HttpPort)

	err := r.Run(listenUrl)
	if err != nil {
		log.Printf("[Error]: Service Start Fail: %s", err.Error())
		panic(fmt.Errorf("[Service]: Service Start Fail: %s", err.Error()))
	}

}
