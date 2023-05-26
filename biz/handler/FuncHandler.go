package handler

import (
	"admin/biz/service"
	"admin/model"
	"admin/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func GetFuncInfo(ctx *gin.Context) {
	allFunc, err := service.GetFunctionList()
	if err != nil {
		errMsg := fmt.Sprintf("[Get FuncInfo] Get FuncInfo Error: %s", err.Error())
		log.Print(errMsg)
		jsonResp := utils.SetBadRequestResp(nil, errMsg)
		ctx.JSON(http.StatusInternalServerError, jsonResp)
		return
	}
	// 返回
	respData := model.GetListResp{
		Total: len(allFunc),
		Items: allFunc,
	}
	jsonResp := utils.SetOKResp(respData, nil)
	ctx.JSON(http.StatusOK, jsonResp)
	return
}

func StartFuncHandler(ctx *gin.Context) {
	functionIDStr := ctx.Param("function_id")
	functionID, err := strconv.ParseInt(functionIDStr, 10, 64)
	if err != nil {
		errMsg := fmt.Sprintf("[Start Function] Parse Request Parameter Error: %s", err.Error())
		log.Print(errMsg)
		jsonResp := utils.SetBadRequestResp(nil, errMsg)
		ctx.JSON(http.StatusBadRequest, jsonResp)
		return
	}

	// 启动镜像
	err = service.StartFunc(functionID)

	if err != nil {
		errMsg := fmt.Sprintf("[Start Function] Start Function Error: %s", err.Error())
		log.Print(errMsg)
		jsonResp := utils.SetBadRequestResp(nil, errMsg)
		ctx.JSON(http.StatusInternalServerError, jsonResp)
		return
	}

	jsonResp := utils.SetOKResp(nil, nil)
	ctx.JSON(http.StatusOK, jsonResp)
	return
}

func StopFuncHandler(ctx *gin.Context) {
	functionIDStr := ctx.Param("function_id")
	functionID, err := strconv.ParseInt(functionIDStr, 10, 64)

	if err != nil {
		errMsg := fmt.Sprintf("[Stop Function] Parse Request Parameter Error: %s", err.Error())
		log.Print(errMsg)
		jsonResp := utils.SetBadRequestResp(nil, errMsg)
		ctx.JSON(http.StatusBadRequest, jsonResp)
		return
	}

	// 停止函数
	err = service.StopFunc(functionID)

	if err != nil {
		errMsg := fmt.Sprintf("[Stop Function] Stop Function Error: %s", err.Error())
		log.Print(errMsg)
		jsonResp := utils.SetBadRequestResp(nil, errMsg)
		ctx.JSON(http.StatusInternalServerError, jsonResp)
		return
	}

	jsonResp := utils.SetOKResp(nil, nil)
	ctx.JSON(http.StatusOK, jsonResp)
	return
}

func DeleteFuncHandler(ctx *gin.Context) {
	functionIDStr := ctx.Param("function_id")
	functionID, err := strconv.ParseInt(functionIDStr, 10, 64)

	if err != nil {
		errMsg := fmt.Sprintf("[Delete Function] Parse Request Parameter Error: %s", err.Error())
		log.Print(errMsg)
		jsonResp := utils.SetBadRequestResp(nil, errMsg)
		ctx.JSON(http.StatusBadRequest, jsonResp)
		return
	}

	// 删除函数
	err = service.DeleteFunc(functionID)

	if err != nil {
		errMsg := fmt.Sprintf("[Delete Function] Delete Function Error: %s", err.Error())
		log.Print(errMsg)
		jsonResp := utils.SetBadRequestResp(nil, errMsg)
		ctx.JSON(http.StatusInternalServerError, jsonResp)
		return
	}
	jsonResp := utils.SetOKResp(nil, nil)
	ctx.JSON(http.StatusOK, jsonResp)
	return
}
