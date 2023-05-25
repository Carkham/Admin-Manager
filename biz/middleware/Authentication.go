package middlerware

import (
	"admin/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const UserCtxKey = "user_info"

type LoginInfo struct {
	ID       int64  `json:"id" `
	Username string `json:"username"`
	IsStaff  int64  `json:"is_staff"`
}

func AuthMiddleware(ctx *gin.Context) {
	sessionID, err := ctx.Cookie("session_id")
	if err != nil {
		ctx.Abort()
		errMsg := fmt.Sprintf("[Authentication] Authenticate Error: %s", err.Error())
		log.Print(errMsg)
		jsonResp := utils.SetUnAuthResp(nil, errMsg)
		ctx.JSON(http.StatusUnauthorized, jsonResp)
		ctx.Abort()
		return
	}
	redisResp := utils.RedisClient.Get(ctx, sessionID)
	if redisResp.Err() != nil {
		ctx.Abort()
		errMsg := fmt.Sprintf("[Authentication] Authenticate Error: %s", redisResp.Err().Error())
		log.Print(errMsg)
		jsonResp := utils.SetUnAuthResp(nil, errMsg)
		ctx.JSON(http.StatusUnauthorized, jsonResp)
		ctx.Abort()
		return
	}

	var loginInfo LoginInfo

	err = json.Unmarshal([]byte(redisResp.Val()), &loginInfo)

	if err != nil {
		ctx.Abort()
		errMsg := fmt.Sprintf("[Authentication] Authenticate Error: %s", err.Error())
		log.Print(errMsg)
		jsonResp := utils.SetUnAuthResp(nil, errMsg)
		ctx.JSON(http.StatusUnauthorized, jsonResp)
		ctx.Abort()
		return
	}

	if loginInfo.IsStaff == 0 {
		ctx.Abort()
		errMsg := fmt.Sprintf("[Authentication] Authenticate Error: Permission Denied")
		log.Print(errMsg)
		jsonResp := utils.SetUnAuthResp(nil, errMsg)
		ctx.JSON(http.StatusUnauthorized, jsonResp)
		ctx.Abort()
		return
	}

	ctx.Set(UserCtxKey, loginInfo)
	ctx.Next()
}
