package handler

import (
	"admin/model"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetNodeInfo(t *testing.T) {
	// 初始化测试数据库
	gin.SetMode(gin.TestMode)
	err := model.TestDBInit()
	if err != nil {
		t.Errorf("init test db error: %v", err)
	}

	// 构造请求
	var testCases = []struct {
		Request            GetReq
		expectedStatusCode int
		expectedResponse   model.JSONResp
	}{
		{
			Request:            GetReq{},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: model.JSONResp{
				Code:    -1,
				Message: "参数格式错误",
				Extra:   "[Node Info] Get Node Info Error: test",
			},
		},
		{
			Request:            GetReq{},
			expectedStatusCode: http.StatusOK,
			expectedResponse: model.JSONResp{
				Message: "success",
				Data: model.GetListResp{
					Total: 1,
					Items: []model.NodeInfo{
						model.NodeInfo{},
					},
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

		GetNodeInfo(c)

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
