package middlerware

import (
	"admin/model"
	"admin/utils"
	"encoding/json"
	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	// 初始化测试数据库
	gin.SetMode(gin.TestMode)
	err := model.TestDBInit()
	if err != nil {
		t.Errorf("InitTestDB error: %s", err.Error())
	}

	// 初始化测试redis
	s := miniredis.RunT(t)
	// set some keys like your codes expected
	loginInfo := `{"id": 1, "username": "fuc_sun", "is_staff": 1}`
	err = s.Set("woshisun", loginInfo)
	if err != nil {
		t.Errorf("set redis error: %s", err.Error())
	}
	utils.InitTestRedisClient(s)
	s.CheckGet(t, "woshisun", loginInfo)

	// 测试样例
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/admin/users", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("Cookie", "session_id=woshisun")
	//c.SetCookie("session_id", "woshisun", 3600, "/", "localhost", false, true)

	AuthMiddleware(c)

	var expected LoginInfo
	err = json.Unmarshal([]byte(loginInfo), &expected)
	actual, _ := c.Get(UserCtxKey)
	assert.Equal(t, expected, actual)
}
