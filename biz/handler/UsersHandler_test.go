package handler

import (
	"admin/model"
	UserUser "admin/model/model"
	ormModel "admin/model/model"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

func TestGetUsersInfo(t *testing.T) {
	// 初始化测试数据库
	gin.SetMode(gin.TestMode)
	err := model.TestDBInit()
	if err != nil {
		t.Errorf("init test db error: %v", err)
	}

	// 数据库预置数据
	_, _ = model.Q.UserUser.Where(model.Q.UserUser.ID.Lte(1000)).Delete()
	testUser := UserUser.UserUser{
		Username: "a",
		Password: "a",
		Email:    "a",
		IsStaff:  true,
	}
	err = model.Q.UserUser.Create(&testUser)
	if err != nil {
		panic(err)
	}

	// 构造正确结果
	var userInfoList []model.UserInfo
	newItem := model.UserInfo{
		UserId:    0,
		UserName:  testUser.Username,
		IsAdmin:   testUser.IsStaff,
		LastLogin: "NONE",
		Email:     testUser.Email,
	}
	userInfoList = append(userInfoList, newItem)

	// 构造请求
	var testCases = []struct {
		Request            GetReq
		expectedStatusCode int
		expectedResponse   model.JSONResp
	}{
		{
			Request:            GetReq{},
			expectedStatusCode: http.StatusOK,
			expectedResponse: model.JSONResp{
				Message: "success",
				Data: model.GetListResp{
					Total: len(userInfoList),
					Items: userInfoList,
				},
			},
		},
	}

	// 测试
	for _, tc := range testCases {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(tc.Request)
		c.Request, _ = http.NewRequest("GET", "/api/admin/users", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		GetUsersInfo(c)

		if w.Code != tc.expectedStatusCode {
			t.Errorf("expected status code %d, got %d", tc.expectedStatusCode, w.Code)
		}

		//print(w.Body)
		var resp model.JSONResp
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		if err != nil {
			t.Errorf("unmarshal response error: %v", err)
		}

		expectedRespJsonBody, _ := json.Marshal(tc.expectedResponse)
		var expObj model.JSONResp
		_ = json.Unmarshal(expectedRespJsonBody, &expObj)
		equal := reflect.DeepEqual(resp, expObj)
		if !equal {
			t.Errorf("response is not expected response %v and %v", resp, expObj)
		}
	}
}

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
		{
			Request:            model.CreateUserReq{Username: "sunshuaibi", Email: "sunshuaibi", IsAdmin: true, Password: "123", NotEncrypt: true},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: model.JSONResp{
				Code:    -1,
				Message: "参数格式错误",
				Extra:   "[Create User] Create User Error: duplicated username or email",
			},
		},
		{
			Request:            model.CreateUserReq{Username: "nn", Email: "nn", IsAdmin: true, Password: "123", NotEncrypt: false},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: model.JSONResp{
				Code:    -1,
				Message: "参数格式错误",
				Extra:   "[Create User] Create User Error: illegal base64 data at input byte 0",
			},
		},
		{
			Request:            model.CreateUserReq{Username: "nn", Email: "nn", IsAdmin: true, Password: "cnVub29i", NotEncrypt: false},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: model.JSONResp{
				Code:    -1,
				Message: "参数格式错误",
				Extra:   "[Create User] Create User Error: open private.pem: The system cannot find the file specified.",
			},
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

		expectedRespJsonBody, _ := json.Marshal(tc.expectedResponse)
		var expObj model.JSONResp
		_ = json.Unmarshal(expectedRespJsonBody, &expObj)
		equal := reflect.DeepEqual(resp, expObj)
		if !equal {
			t.Errorf("response is not expected response %v and %v", resp, expObj)
		}
	}

	var testCases2 = []struct {
		Request            string
		expectedStatusCode int
		expectedResponse   model.JSONResp
	}{
		{
			Request:            "sss",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: model.JSONResp{
				Code:    -1,
				Message: "参数格式错误",
				Extra:   "[Create User] Parse Parameter Error: json: cannot unmarshal string into Go value of type model.CreateUserReq",
			},
		},
	}

	for _, tc := range testCases2 {
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

type GetReq struct {
}

func TestDeleteUser(t *testing.T) {
	// 初始化测试数据库
	gin.SetMode(gin.TestMode)
	err := model.TestDBInit()
	if err != nil {
		t.Errorf("init test db error: %v", err)
	}

	// 数据库预置数据
	testUser := UserUser.UserUser{
		ID:       10,
		Username: "a",
		Password: "a",
		Email:    "a",
		IsStaff:  true,
	}
	err = model.Q.UserUser.Create(&testUser)
	if err != nil {
		panic(err)
	}
	// 构造测试数据
	newFunc := ormModel.Function{
		FunctionID:    10,
		FunctionLabel: "10",
		UserID:        10,
		TriggerID:     10,
		SrcType:       "10",
		SrcLoc:        "10",
		Replicas:      10,
		QuotaInfo:     "10",
		TemplateID:    10,
	}
	err = model.Q.Function.Create(&newFunc)
	if err != nil {
		panic(err)
	}

	var testCases = []struct {
		Request            GetReq
		expectedStatusCode int
		expectedResponse   model.JSONResp
	}{
		{
			Request:            GetReq{},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: model.JSONResp{
				Code:    -1,
				Message: "参数格式错误",
				Extra:   "[Delete User] Delete User Error: function is still running",
			},
		},
		{
			Request:            GetReq{},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   model.JSONResp{Message: "success"},
		},
	}

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(tc.Request)
		c.Request, _ = http.NewRequest("DELETE", "/api/admin/user/0", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.AddParam("user_id", strconv.FormatInt(testUser.ID, 10))

		DeleteUser(c)

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

	var testCases2 = []struct {
		Request            string
		expectedStatusCode int
		expectedResponse   model.JSONResp
	}{
		{
			Request:            "sss",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: model.JSONResp{
				Code:    -1,
				Message: "参数格式错误",
				Extra:   "[Delete User] Parse Parameter Error: strconv.ParseInt: parsing \"s\": invalid syntax",
			},
		},
	}

	for _, tc := range testCases2 {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(tc.Request)
		c.Request, _ = http.NewRequest("PUT", "/api/admin/users", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.AddParam("user_id", "s")

		DeleteUser(c)

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
