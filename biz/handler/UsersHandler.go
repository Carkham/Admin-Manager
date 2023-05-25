package handler

import (
	"admin/model"
	UserUser "admin/model/model"
	"admin/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetUsersInfo(ctx *gin.Context) {
	var userUserList []UserUser.UserUser
	err := model.Q.UserUser.Scan(&userUserList)
	if err != nil {
		errMsg := fmt.Sprintf("[Get UsersInfo] Get UsersInfo Error: %s", err.Error())
		log.Print(errMsg)
		jsonResp := utils.SetBadRequestResp(nil, errMsg)
		ctx.JSON(http.StatusInternalServerError, jsonResp)
		return
	}

	var userInfoList []model.UserInfo
	for _, userUser := range userUserList {
		userInfoList = append(userInfoList, model.UserInfo{
			UserId:    int(userUser.ID),
			UserName:  userUser.Username,
			LastLogin: userUser.LastLogin.Format("YYYY-MM-DD hh:mm"),
			IsAdmin:   userUser.IsStaff,
			Email:     userUser.Email,
		})
	}

	// 返回
	respData := model.GetListResp{
		Total: len(userInfoList),
		Items: userInfoList,
	}
	jsonResp := utils.SetOKResp(respData, nil)
	ctx.JSON(http.StatusOK, jsonResp)
	return
}
