package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserHandler_Update(t *testing.T) {
	r, h := setupUserRouter()
	r.PUT("/api/system/user", h.Update)

	body := `{"id":1001,"nickname":"更新后","email":"new@test.com"}`
	req := httptest.NewRequest(http.MethodPut, "/api/system/user", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])
	assert.Equal(t, "更新成功", resp["msg"])
}

func TestUserHandler_Delete(t *testing.T) {
	r, h := setupUserRouter()
	r.DELETE("/api/system/user/:id", h.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/api/system/user/1001", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])
}

func TestUserHandler_Delete_InvalidID(t *testing.T) {
	r, h := setupUserRouter()
	r.DELETE("/api/system/user/:id", h.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/api/system/user/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(400), resp["code"])
}

func TestUserHandler_Page(t *testing.T) {
	r, h := setupUserRouter()
	r.GET("/api/system/user/page", h.Page)

	req := httptest.NewRequest(http.MethodGet, "/api/system/user/page?page=1&pageSize=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])
}

func TestUserHandler_Get_InvalidID(t *testing.T) {
	r, h := setupUserRouter()
	r.GET("/api/system/user/:id", h.Get)

	req := httptest.NewRequest(http.MethodGet, "/api/system/user/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(400), resp["code"])
}
