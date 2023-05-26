package handler

import (
	"admin/model"
	"admin/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetNodeInfo(ctx *gin.Context) {
	nodesMap, err := utils.GetNodeList()
	nodesList := make([]model.NodeInfo, 0, len(nodesMap))
	for _, v := range nodesMap {
		nodesList = append(nodesList, *v)
	}

	if err != nil {
		errMsg := fmt.Sprintf("[Node Info] Get Node Info Error: %s", err.Error())
		log.Print(errMsg)
		jsonResp := utils.SetBadRequestResp(nil, errMsg)
		ctx.JSON(http.StatusInternalServerError, jsonResp)
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
