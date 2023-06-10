package handler

import (
	"admin/model"
	Template "admin/model/model"
	"admin/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateTemplate(ctx *gin.Context) {
	var req model.CreateTemplateReq
	err := ctx.Bind(&req)
	if err != nil {
		throwError(ctx, fmt.Sprintf("[Create Template] Parse Parameter Error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// 关键参数不能为空
	if req.TemplateLabel == "" || req.ImageName == "" {
		throwError(ctx, fmt.Sprintf("[Create Template] Parse Parameter Error: Parameter can not be none"), http.StatusBadRequest)
		return
	}

	// 工程模板文件名可以为空
	if req.BaseCode == "" {
		req.BaseCode = req.TemplateLabel + ".zip"
	}

	var newTemplate = Template.Template{
		ImageName:     req.ImageName,
		TemplateLabel: req.TemplateLabel,
		FileName:      req.BaseCode,
	}
	err = model.Q.Template.Create(&newTemplate)
	if err != nil {
		throwError(ctx, fmt.Sprintf("[Create Template] Create Template Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// 返回
	respData := model.CreateTemplateResp{
		TemplateId: int(newTemplate.TemplateID),
	}
	jsonResp := utils.SetOKResp(respData, nil)
	ctx.JSON(http.StatusOK, jsonResp)
	return
}
