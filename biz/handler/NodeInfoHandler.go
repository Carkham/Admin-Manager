package handler

import (
	"admin/model"
	"admin/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetNodeInfo(ctx *gin.Context) {
	var nodesMap map[string]*model.NodeInfo
	var err error
	nodesMap, err = utils.GetNodeList()

	nodesList := make([]model.NodeInfo, 0, len(nodesMap))
	for _, v := range nodesMap {
		nodesList = append(nodesList, *v)
	}

	if err != nil {
		throwError(ctx, fmt.Sprintf("[Node Info] Get Node Info Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	respData := model.GetListResp{
		Total: len(nodesList),
		Items: nodesList,
	}
	jsonResp := utils.SetOKResp(respData, nil)
	ctx.JSON(http.StatusOK, jsonResp)
	return
}
