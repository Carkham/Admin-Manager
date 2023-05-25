package handler

import (
	"admin/model"
	"admin/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func CreateTemplate(ctx *gin.Context) {
	var req model.CreateTemplateReq
	err := ctx.Bind(&req)
	if err != nil {
		errMsg := fmt.Sprintf("[Create Template] Parse Parameter Error: %s", err.Error())
		log.Print(errMsg)
		jsonResp := utils.SetBadRequestResp(nil, errMsg)
		ctx.JSON(http.StatusBadRequest, jsonResp)
		return
	}

	// 关键参数不能为空
	if req.TemplateLabel == "" || req.ImageName == "" {
		errMsg := fmt.Sprintf("[Create Template] Parse Parameter Error: Parameter can not be none")
		log.Print(errMsg)
		jsonResp := utils.SetBadRequestResp(nil, errMsg)
		ctx.JSON(http.StatusBadRequest, jsonResp)
		return
	}

	// 镜像名称可以为空
	if req.BaseCode == "" {
		req.BaseCode = req.TemplateLabel + ".zip"
	}

	// todo 创建镜像

	// 返回
	respData := model.CreateTemplateResp{
		TemplateId: 0,
	}
	jsonResp := utils.SetOKResp(respData, nil)
	ctx.JSON(http.StatusOK, jsonResp)
	return
}
