package handler

import (
	"admin/biz/service"
	"admin/model"
	Function "admin/model/model"
	UserUser "admin/model/model"
	"admin/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
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
			LastLogin: userUser.LastLogin.Format("2006-01-02 15:04"),
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

func CreateUser(ctx *gin.Context) {
	var req model.CreateUserReq
	err := ctx.Bind(&req)
	if err != nil {
		errMsg := fmt.Sprintf("[Create User] Parse Parameter Error: %s", err.Error())
		log.Print(errMsg)
		jsonResp := utils.SetBadRequestResp(nil, errMsg)
		ctx.JSON(http.StatusBadRequest, jsonResp)
		return
	}

	// 用户名密码不能为空
	if req.Username == "" || req.Password == "" {
		errMsg := fmt.Sprintf("[Create User] Parse Parameter Error: need username or password")
		log.Print(errMsg)
		jsonResp := utils.SetBadRequestResp(nil, errMsg)
		ctx.JSON(http.StatusBadRequest, jsonResp)
		return
	}

	// 防止邮箱用户名重复
	var duplicatedUsers []UserUser.UserUser
	err = model.Q.UserUser.Where(
		model.Q.UserUser.Username.Eq(req.Username),
	).Or(
		model.Q.UserUser.Email.Eq(req.Email),
	).Scan(&duplicatedUsers)
	if err != nil {
		errMsg := fmt.Sprintf("[Create User] Create User Error: %s", err.Error())
		log.Print(errMsg)
		jsonResp := utils.SetBadRequestResp(nil, errMsg)
		ctx.JSON(http.StatusInternalServerError, jsonResp)
		return
	}
	if len(duplicatedUsers) > 0 {
		errMsg := fmt.Sprintf("[Create User] Create User Error: duplicated username or email")
		log.Print(errMsg)
		jsonResp := utils.SetBadRequestResp(nil, errMsg)
		ctx.JSON(http.StatusBadRequest, jsonResp)
		return
	}

	if req.NotEncrypt != true {
		// rsa加密密码解密
		cipherText, err := utils.DecodeBase64(req.Password)
		if err != nil {
			errMsg := fmt.Sprintf("[Create User] Create User Error: %s", err.Error())
			log.Print(errMsg)
			jsonResp := utils.SetBadRequestResp(nil, errMsg)
			ctx.JSON(http.StatusInternalServerError, jsonResp)
			return
		}
		req.Password, err = utils.RSADecrypt(cipherText)
		if err != nil {
			errMsg := fmt.Sprintf("[Create User] Create User Error: %s", err.Error())
			log.Print(errMsg)
			jsonResp := utils.SetBadRequestResp(nil, errMsg)
			ctx.JSON(http.StatusInternalServerError, jsonResp)
			return
		}
	}

	// 密码存入数据库前加密
	encodeParam := utils.EncodeParam{
		Memory:      102400,
		Iterations:  2,
		Parallelism: 8,
		SaltLength:  22,
		KeyLength:   32,
	}

	hash, err := utils.GenerateFromPassword(req.Password, encodeParam)
	if err != nil {
		errMsg := fmt.Sprintf("[Create User] Create User Error: %s", err.Error())
		log.Print(errMsg)
		jsonResp := utils.SetBadRequestResp(nil, errMsg)
		ctx.JSON(http.StatusInternalServerError, jsonResp)
		return
	}

	// 保存
	var newUser = UserUser.UserUser{
		Username: req.Username,
		Email:    req.Email,
		Password: hash,
		IsStaff:  req.IsAdmin,
	}
	err = model.Q.UserUser.Create(&newUser)
	if err != nil {
		errMsg := fmt.Sprintf("[Create User] Create User Error: %s", err.Error())
		log.Print(errMsg)
		jsonResp := utils.SetBadRequestResp(nil, errMsg)
		ctx.JSON(http.StatusInternalServerError, jsonResp)
		return
	}

	// 完成返回
	respData := model.CreateUserResp{
		UserId: newUser.ID,
	}
	jsonResp := utils.SetOKResp(respData, nil)
	ctx.JSON(http.StatusOK, jsonResp)
	return
}

func DeleteUser(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		errMsg := fmt.Sprintf("[Delete User] Parse Parameter Error: %s", err.Error())
		log.Print(errMsg)
		jsonResp := utils.SetBadRequestResp(nil, errMsg)
		ctx.JSON(http.StatusBadRequest, jsonResp)
		return
	}

	// 检查函数是否还在运行
	var functions []Function.Function
	err = model.Q.Function.Where(model.Q.Function.UserID.Eq(userID)).Scan(&functions)
	for _, v := range functions {
		runningList, _ := utils.GetPodInfoList(v.FunctionID)
		if len(runningList) > 0 {
			errMsg := fmt.Sprintf("[Delete User] Delete User Error: function is still running")
			log.Print(errMsg)
			jsonResp := utils.SetServerErrorResp(nil, errMsg)
			ctx.JSON(http.StatusBadRequest, jsonResp)
			return
		}
	}

	// 删除所有函数
	for _, v := range functions {
		_ = service.DeleteFunc(v.FunctionID)
	}

	_, err = model.Q.UserUser.Where(model.Q.UserUser.ID.Eq(userID)).Delete()
	if err != nil {
		errMsg := fmt.Sprintf("[Delete User] Delete User Error: %s", err.Error())
		log.Print(errMsg)
		jsonResp := utils.SetBadRequestResp(nil, errMsg)
		ctx.JSON(http.StatusInternalServerError, jsonResp)
		return
	}

	jsonResp := utils.SetOKResp(nil, nil)
	ctx.JSON(http.StatusOK, jsonResp)
	return
}
