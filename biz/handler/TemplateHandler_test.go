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

func TestCreateHandlerInfo(t *testing.T) {
	// 初始化测试数据库
	gin.SetMode(gin.TestMode)
	err := model.TestDBInit()
	if err != nil {
		t.Errorf("init test db error: %v", err)
	}

	// 构造请求
	var testCases = []struct {
		Request            model.CreateTemplateReq
		expectedStatusCode int
		expectedResponse   model.JSONResp
	}{
		{
			Request: model.CreateTemplateReq{
				TemplateLabel: "",
				ImageName:     "1",
				BaseCode:      "1",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: model.JSONResp{
				Code:    -1,
				Message: "参数格式错误",
				Extra:   "[Create Template] Parse Parameter Error: Parameter can not be none",
			},
		},
		{
			Request: model.CreateTemplateReq{
				TemplateLabel: "1",
				ImageName:     "1",
				BaseCode:      "",
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: model.JSONResp{
				Code:    0,
				Message: "success",
				Data: model.CreateTemplateResp{
					TemplateId: 0,
				},
			},
		},
	}

	// 测试
	for _, tc := range testCases {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(tc.Request)
		c.Request, _ = http.NewRequest("PUT", "/api/admin/template", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		CreateTemplate(c)

		if w.Code != tc.expectedStatusCode {
			t.Errorf("expected status code %d, got %d", tc.expectedStatusCode, w.Code)
		}

		//print(w.Body)
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
			t.Errorf("expected response %s, got %s", expectedString, respString)
		}
	}

	// 构造请求
	var testCases2 = []struct {
		Request            string
		expectedStatusCode int
		expectedResponse   model.JSONResp
	}{
		{
			Request:            "Sss",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: model.JSONResp{
				Code:    -1,
				Message: "参数格式错误",
				Extra:   "[Create Template] Parse Parameter Error: json: cannot unmarshal string into Go value of type model.CreateTemplateReq",
			},
		},
	}

	// 测试
	for _, tc := range testCases2 {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(tc.Request)
		c.Request, _ = http.NewRequest("PUT", "/api/admin/template", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		CreateTemplate(c)

		if w.Code != tc.expectedStatusCode {
			t.Errorf("expected status code %d, got %d", tc.expectedStatusCode, w.Code)
		}

		//print(w.Body)
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
			t.Errorf("expected response %s, got %s", expectedString, respString)
		}
	}
}
