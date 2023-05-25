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
