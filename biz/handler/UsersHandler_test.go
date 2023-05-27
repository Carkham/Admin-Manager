package handler

import (
	"admin/model"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUser(t *testing.T) {
	// 初始化测试数据库
	gin.SetMode(gin.TestMode)
	err := model.TestDBInit()
	if err != nil {
		t.Errorf("init test db error: %v", err)
	}

	var testCases = []struct {
		Request            model.CreateUserReq
		expectedStatusCode int
		expectedResponse   model.JSONResp
	}{
		{
			Request:            model.CreateUserReq{Username: "sunshuaibi", Email: "sunshuaibi", IsAdmin: true, Password: "123", NotEncrypt: true},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   model.JSONResp{Message: "success", Data: model.CreateUserResp{UserId: 0}},
		},
		{
			Request:            model.CreateUserReq{Username: "sun_shuaibi", Email: "sun_shuaibi", NotEncrypt: true},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   model.JSONResp{Code: -1, Message: "参数格式错误", Extra: "[Create User] Parse Parameter Error: need username or password"},
		},
	}

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(tc.Request)
		c.Request, _ = http.NewRequest("PUT", "/api/admin/users", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		CreateUser(c)

		if w.Code != tc.expectedStatusCode {
			t.Errorf("expected status code %d, got %d", tc.expectedStatusCode, w.Code)
		}

		var resp model.JSONResp
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		if err != nil {
			t.Errorf("unmarshal response error: %v", err)
		}

		respJsonBody, _ := json.Marshal(resp)
		expectedRespJsonBody, _ := json.Marshal(tc.expectedResponse)
		respString := string(respJsonBody)
		expectedString := string(expectedRespJsonBody)

		if respString != expectedString {
			t.Errorf("expected response %v, got %v", expectedString, respString)
		}
	}
}
