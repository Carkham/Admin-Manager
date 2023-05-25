package handler

import (
	"admin/model"
	"admin/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetNodeInfo(ctx *gin.Context) {
	var nodes []model.NodeItem
	// todo 获取节点信息
	respData := model.GetListResp{
		Total: len(nodes),
		Items: nodes,
	}
	jsonResp := utils.SetOKResp(respData, nil)
	ctx.JSON(http.StatusOK, jsonResp)
	return
}
