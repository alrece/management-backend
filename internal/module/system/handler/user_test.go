package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	smodel "management-backend/internal/module/system/model"
	"management-backend/internal/module/system/service"
	"management-backend/testutil"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// mockUserService 测试用 mock Service
type mockUserService struct {
	users map[int64]*smodel.UserResp
}

func newMockUserService() *mockUserService {
	return &mockUserService{users: make(map[int64]*smodel.UserResp)}
}

func (m *mockUserService) Create(_ context.Context, req *smodel.UserCreateReq, creator int64, tenantID int64) (int64, error) {
	return 9001, nil
}

func (m *mockUserService) Update(_ context.Context, req *smodel.UserUpdateReq, updater int64) error {
	return nil
}

func (m *mockUserService) Delete(_ context.Context, id int64) error {
	return nil
}

func (m *mockUserService) GetByID(_ context.Context, id int64) (*smodel.UserResp, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserService) GetByUsername(_ context.Context, username string) (*smodel.User, error) {
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserService) GetRawByID(_ context.Context, id int64) (*smodel.User, error) {
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserService) UpdateProfile(_ context.Context, req *service.UserProfileUpdateReq) error {
	return nil
}

func (m *mockUserService) UpdatePassword(_ context.Context, userID int64, hashedPassword string) error {
	return nil
}

func (m *mockUserService) Page(_ context.Context, req *smodel.UserPageReq) ([]smodel.UserResp, int64, error) {
	return []smodel.UserResp{}, 0, nil
}

func setupUserRouter() (*gin.Engine, *UserHandler) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := newMockUserService()
	h := NewUserHandler(svc)
	return r, h
}

func TestUserHandler_Create(t *testing.T) {
	testutil.NewTestDB(t)
	r, h := setupUserRouter()
	r.POST("/api/system/user", h.Create)

	body := `{"username":"newuser","password":"pass123","nickname":"新用户"}`
	req := httptest.NewRequest(http.MethodPost, "/api/system/user", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])
}

func TestUserHandler_Create_BindError(t *testing.T) {
	r, h := setupUserRouter()
	r.POST("/api/system/user", h.Create)

	req := httptest.NewRequest(http.MethodPost, "/api/system/user", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.NotEqual(t, float64(0), resp["code"])
}

func TestUserHandler_Get(t *testing.T) {
	r, h := setupUserRouter()
	r.GET("/api/system/user/:id", h.Get)

	mockSvc := h.svc.(*mockUserService)
	mockSvc.users[1001] = &smodel.UserResp{
		ID:       1001,
		Username: "testuser",
		Nickname: "测试",
	}

	req := httptest.NewRequest(http.MethodGet, "/api/system/user/1001", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])
}
