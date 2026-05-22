package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestOk(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	Ok(c, map[string]string{"name": "test"})

	assert.Equal(t, http.StatusOK, w.Code)
	var r R
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &r))
	assert.Equal(t, 0, r.Code)
	assert.Equal(t, "success", r.Message)
}

func TestOkMsg(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	OkMsg(c, "创建成功")

	assert.Equal(t, http.StatusOK, w.Code)
	var r R
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &r))
	assert.Equal(t, 0, r.Code)
	assert.Equal(t, "创建成功", r.Message)
}

func TestOkPage(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	list := []map[string]string{{"name": "a"}, {"name": "b"}}
	OkPage(c, list, 100, 1, 10)

	assert.Equal(t, http.StatusOK, w.Code)
	var r struct {
		Code int        `json:"code"`
		Msg  string     `json:"msg"`
		Data PageResult `json:"data"`
	}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &r))
	assert.Equal(t, 0, r.Code)
	assert.Equal(t, int64(100), r.Data.Total)
	assert.Equal(t, 1, r.Data.Page)
	assert.Equal(t, 10, r.Data.PageSize)
}

func TestFail(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	Fail(c, 1001, "用户不存在")

	assert.Equal(t, http.StatusOK, w.Code)
	var r R
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &r))
	assert.Equal(t, 1001, r.Code)
	assert.Equal(t, "用户不存在", r.Message)
}

func TestFailWithStatus(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	FailWithStatus(c, http.StatusUnauthorized, 401, "未授权")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var r R
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &r))
	assert.Equal(t, 401, r.Code)
	assert.Equal(t, "未授权", r.Message)
}
