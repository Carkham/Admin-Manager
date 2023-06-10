package handler

import (
	"admin/model"
	ormModel "admin/model/model"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetFuncInfo(t *testing.T) {
	// 初始化测试数据库
	gin.SetMode(gin.TestMode)
	err := model.TestDBInit()
	if err != nil {
		t.Errorf("init test db error: %v", err)
	}

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
					Total: 0,
					Items: nil,
				},
			},
		},
		{
			Request:            GetReq{},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: model.JSONResp{
				Code:    -1,
				Message: "参数格式错误",
				Extra:   "[Get FuncInfo] Get FuncInfo Error: test",
			},
		},
	}

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(tc.Request)
		c.Request, _ = http.NewRequest("GET", "/api/admin/functions", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		GetFuncInfo(c)

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
}

func TestStartFuncHandler(t *testing.T) {
	// 初始化测试数据库
	gin.SetMode(gin.TestMode)
	err := model.TestDBInit()
	if err != nil {
		t.Errorf("init test db error: %v", err)
	}

	// 构造测试数据
	newFunc := ormModel.Function{
		FunctionID:    100,
		FunctionLabel: "100",
		UserID:        100,
		TriggerID:     100,
		SrcType:       "100",
		SrcLoc:        "100",
		Replicas:      100,
		QuotaInfo:     "100",
		TemplateID:    100,
	}
	err = model.Q.Function.Create(&newFunc)
	if err != nil {
		panic(err)
	}
	triggerLabelTest := "100"
	var userIdTest int64
	userIdTest = 100
	newTrigger := ormModel.Trigger{
		TriggerID:     100,
		TriggerType:   "100",
		TriggerConfig: "100",
		TriggerLabel:  &triggerLabelTest,
		UserID:        &userIdTest,
	}
	err = model.Q.Trigger.Create(&newTrigger)
	if err != nil {
		panic(err)
	}
	newTemplate := ormModel.Template{
		TemplateID:    100,
		ImageName:     "100",
		TemplateLabel: "100",
		FileName:      "100",
	}
	err = model.Q.Template.Create(&newTemplate)
	if err != nil {
		panic(err)
	}
	newUser := ormModel.UserUser{
		ID:       100,
		Username: "100",
		Password: "100",
		Email:    "100",
		IsStaff:  true,
	}
	_ = model.Q.UserUser.Create(&newUser)

	var testCases = []struct {
		Param              string
		Request            GetReq
		expectedStatusCode int
		expectedResponse   model.JSONResp
	}{
		{
			Param:              "100",
			Request:            GetReq{},
			expectedStatusCode: http.StatusOK,
			expectedResponse: model.JSONResp{
				Message: "success",
			},
		},
		{
			Param:              "0",
			Request:            GetReq{},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: model.JSONResp{
				Code:    -1,
				Message: "参数格式错误",
				Extra:   "[Start Function] Start Function Error: no such function 0",
			},
		},
		{
			Param:              "g",
			Request:            GetReq{},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: model.JSONResp{
				Code:    -1,
				Message: "参数格式错误",
				Extra:   "[Start Function] Parse Request Parameter Error: strconv.ParseInt: parsing \"g\": invalid syntax",
			},
		},
	}

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(tc.Request)
		c.Request, _ = http.NewRequest("POST", "/api/admin/functions/:function_id/start", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.AddParam("function_id", tc.Param)

		StartFuncHandler(c)

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
}

func TestStopFuncHandler(t *testing.T) {
	// 初始化测试数据库
	gin.SetMode(gin.TestMode)
	err := model.TestDBInit()
	if err != nil {
		t.Errorf("init test db error: %v", err)
	}

	// 构造测试数据
	newFunc := ormModel.Function{
		FunctionID:    101,
		FunctionLabel: "101",
		UserID:        101,
		TriggerID:     101,
		SrcType:       "101",
		SrcLoc:        "101",
		Replicas:      101,
		QuotaInfo:     "101",
		TemplateID:    101,
	}
	err = model.Q.Function.Create(&newFunc)
	if err != nil {
		panic(err)
	}
	triggerLabelTest := "101"
	var userIdTest int64
	userIdTest = 101
	newTrigger := ormModel.Trigger{
		TriggerID:     101,
		TriggerType:   "101",
		TriggerConfig: "101",
		TriggerLabel:  &triggerLabelTest,
		UserID:        &userIdTest,
	}
	err = model.Q.Trigger.Create(&newTrigger)
	if err != nil {
		panic(err)
	}
	newTemplate := ormModel.Template{
		TemplateID:    101,
		ImageName:     "101",
		TemplateLabel: "101",
		FileName:      "101",
	}
	err = model.Q.Template.Create(&newTemplate)
	if err != nil {
		panic(err)
	}
	newUser := ormModel.UserUser{
		ID:       101,
		Username: "101",
		Password: "101",
		Email:    "101",
		IsStaff:  true,
	}
	_ = model.Q.UserUser.Create(&newUser)

	var testCases = []struct {
		Param              string
		Request            GetReq
		expectedStatusCode int
		expectedResponse   model.JSONResp
	}{
		{
			Param:              "101",
			Request:            GetReq{},
			expectedStatusCode: http.StatusOK,
			expectedResponse: model.JSONResp{
				Message: "success",
			},
		},
		{
			Param:              "2",
			Request:            GetReq{},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: model.JSONResp{
				Code:    -1,
				Message: "参数格式错误",
				Extra:   "[Stop Function] Stop Function Error: record not found",
			},
		},
		{
			Param:              "g",
			Request:            GetReq{},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: model.JSONResp{
				Code:    -1,
				Message: "参数格式错误",
				Extra:   "[Stop Function] Parse Request Parameter Error: strconv.ParseInt: parsing \"g\": invalid syntax",
			},
		},
	}

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(tc.Request)
		c.Request, _ = http.NewRequest("POST", "/api/admin/functions/:function_id/stop", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.AddParam("function_id", tc.Param)

		StopFuncHandler(c)

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
}

func TestDeleteFuncHandler(t *testing.T) {
	// 初始化测试数据库
	gin.SetMode(gin.TestMode)
	err := model.TestDBInit()
	if err != nil {
		t.Errorf("init test db error: %v", err)
	}

	// 构造测试数据
	newFunc := ormModel.Function{
		FunctionID:    102,
		FunctionLabel: "102",
		UserID:        102,
		TriggerID:     102,
		SrcType:       "102",
		SrcLoc:        "102",
		Replicas:      102,
		QuotaInfo:     "102",
		TemplateID:    102,
	}
	err = model.Q.Function.Create(&newFunc)
	if err != nil {
		panic(err)
	}
	triggerLabelTest := "102"
	var userIdTest int64
	userIdTest = 102
	newTrigger := ormModel.Trigger{
		TriggerID:     102,
		TriggerType:   "102",
		TriggerConfig: "102",
		TriggerLabel:  &triggerLabelTest,
		UserID:        &userIdTest,
	}
	err = model.Q.Trigger.Create(&newTrigger)
	if err != nil {
		panic(err)
	}
	newTemplate := ormModel.Template{
		TemplateID:    102,
		ImageName:     "102",
		TemplateLabel: "102",
		FileName:      "102",
	}
	err = model.Q.Template.Create(&newTemplate)
	if err != nil {
		panic(err)
	}
	newUser := ormModel.UserUser{
		ID:       102,
		Username: "102",
		Password: "102",
		Email:    "102",
		IsStaff:  true,
	}
	_ = model.Q.UserUser.Create(&newUser)

	var testCases = []struct {
		Param              string
		Request            GetReq
		expectedStatusCode int
		expectedResponse   model.JSONResp
	}{
		{
			Param:              "102",
			Request:            GetReq{},
			expectedStatusCode: http.StatusOK,
			expectedResponse: model.JSONResp{
				Message: "success",
			},
		},
		{
			Param:              "2",
			Request:            GetReq{},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: model.JSONResp{
				Code:    -1,
				Message: "参数格式错误",
				Extra:   "[Delete Function] Delete Function Error: record not found",
			},
		},
		{
			Param:              "g",
			Request:            GetReq{},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: model.JSONResp{
				Code:    -1,
				Message: "参数格式错误",
				Extra:   "[Delete Function] Parse Request Parameter Error: strconv.ParseInt: parsing \"g\": invalid syntax",
			},
		},
	}

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(tc.Request)
		c.Request, _ = http.NewRequest("POST", "/api/admin/functions/:function_id/stop", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.AddParam("function_id", tc.Param)

		DeleteFuncHandler(c)

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
}
