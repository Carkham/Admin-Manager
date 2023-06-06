package main

import (
	"admin/biz/handler"
	"admin/biz/middleware"
	"github.com/gin-gonic/gin"
)

func customRouter(r *gin.Engine) {
	r.Use(middlerware.AuthMiddleware)
	r.GET("/api/admin/nodes", handler.GetNodeInfo)
	r.GET("/api/admin/functions", handler.GetFuncInfo)
	r.PUT("/api/admin/template", handler.CreateTemplate)
	r.POST("/api/admin/functions/:function_id/start", handler.StartFuncHandler)
	r.POST("/api/admin/functions/:function_id/stop", handler.StopFuncHandler)
	r.POST("/api/admin/functions/:function_id", handler.DeleteFuncHandler)
	r.GET("/api/admin/users", handler.GetUsersInfo)
	r.PUT("/api/admin/users", handler.CreateUser)
	r.DELETE("/api/admin/user/:user_id", handler.DeleteUser)
}
