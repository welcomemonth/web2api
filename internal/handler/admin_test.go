package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/welcomemonth/web2api/internal/auth"
	"github.com/welcomemonth/web2api/internal/repository"
	"github.com/welcomemonth/web2api/internal/service"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newLoginRouter(t *testing.T) (*gin.Engine, *auth.JWTService) {
	t.Helper()

	repo := repository.NewAdminRepository(filepath.Join(t.TempDir(), "admin.json"))
	if _, err := service.NewAdminInitService(repo).Ensure(); err != nil {
		t.Fatalf("初始化管理员失败: %v", err)
	}
	jwtSvc := auth.NewJWTService("test-secret", time.Hour)

	r := gin.New()
	NewAdminHandler(service.NewAdminAuthService(repo, jwtSvc)).RegisterAdminRoutes(r)
	return r, jwtSvc
}

func postLogin(r *gin.Engine, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/admin/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestLoginSuccess(t *testing.T) {
	r, jwtSvc := newLoginRouter(t)

	w := postLogin(r, `{"username":"admin","password":"admin123"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("code = %d, want 200, body=%s", w.Code, w.Body.String())
	}

	var resp struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("响应不是合法 JSON: %v", err)
	}
	if resp.Token == "" {
		t.Fatalf("响应缺少 token: %s", w.Body.String())
	}

	claims, err := jwtSvc.Parse(resp.Token)
	if err != nil {
		t.Fatalf("返回的 Token 应可校验: %v", err)
	}
	if claims.Username != "admin" {
		t.Errorf("Username = %q, want admin", claims.Username)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	r, _ := newLoginRouter(t)

	w := postLogin(r, `{"username":"admin","password":"nope"}`)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("错误密码应为 401, got %d", w.Code)
	}
}

func TestLoginMissingFields(t *testing.T) {
	r, _ := newLoginRouter(t)

	cases := []string{
		`{"username":"admin"}`,
		`{"password":"admin123"}`,
		`{}`,
		`{"username":"","password":""}`,
	}
	for _, body := range cases {
		if w := postLogin(r, body); w.Code != http.StatusBadRequest {
			t.Errorf("body=%s 应为 400, got %d", body, w.Code)
		}
	}
}

func TestLoginMalformedJSON(t *testing.T) {
	r, _ := newLoginRouter(t)

	if w := postLogin(r, `{not-json`); w.Code != http.StatusBadRequest {
		t.Errorf("非法 JSON 应为 400, got %d", w.Code)
	}
}
