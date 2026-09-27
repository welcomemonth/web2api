package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/welcomemonth/web2api/internal/auth"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// newTestRouter 构造一个只挂 AdminJWT 的路由，命中即 200 并回显 Context 中的用户名。
func newTestRouter(svc *auth.JWTService) *gin.Engine {
	r := gin.New()
	r.GET("/protected", AdminJWT(svc), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"username": c.GetString(ContextUsername)})
	})
	return r
}

func doRequest(r *gin.Engine, authHeader string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestAdminJWTValidToken(t *testing.T) {
	svc := auth.NewJWTService("secret", time.Hour)
	token, err := svc.Generate("admin")
	if err != nil {
		t.Fatal(err)
	}

	w := doRequest(newTestRouter(svc), "Bearer "+token)
	if w.Code != http.StatusOK {
		t.Fatalf("code = %d, want 200, body=%s", w.Code, w.Body.String())
	}
	if got := w.Body.String(); !strings.Contains(got, "admin") {
		t.Errorf("用户名未写入 Context: %s", got)
	}
}

func TestAdminJWTMissingHeader(t *testing.T) {
	svc := auth.NewJWTService("secret", time.Hour)

	if w := doRequest(newTestRouter(svc), ""); w.Code != http.StatusUnauthorized {
		t.Errorf("无 Authorization 头应为 401, got %d", w.Code)
	}
}

func TestAdminJWTMalformedHeader(t *testing.T) {
	svc := auth.NewJWTService("secret", time.Hour)
	token, _ := svc.Generate("admin")

	cases := []string{
		token,             // 缺少 Bearer 前缀
		"Basic " + token,  // 错误 scheme
		"Bearer",          // 只有 scheme
		"Bearer ",         // 前缀后为空
		"bearer " + token, // 大小写不符（第一版只认 Bearer）
	}

	for _, h := range cases {
		if w := doRequest(newTestRouter(svc), h); w.Code != http.StatusUnauthorized {
			t.Errorf("Authorization=%q 应为 401, got %d", h, w.Code)
		}
	}
}

func TestAdminJWTInvalidToken(t *testing.T) {
	svc := auth.NewJWTService("secret", time.Hour)

	for _, tok := range []string{"garbage", "a.b.c"} {
		if w := doRequest(newTestRouter(svc), "Bearer "+tok); w.Code != http.StatusUnauthorized {
			t.Errorf("Token=%q 应为 401, got %d", tok, w.Code)
		}
	}
}

func TestAdminJWTExpiredToken(t *testing.T) {
	// 用负有效期签发 → 立即过期
	svc := auth.NewJWTService("secret", -time.Hour)
	token, err := svc.Generate("admin")
	if err != nil {
		t.Fatal(err)
	}

	// 校验方用正常服务，密钥一致
	verifier := auth.NewJWTService("secret", time.Hour)
	if w := doRequest(newTestRouter(verifier), "Bearer "+token); w.Code != http.StatusUnauthorized {
		t.Errorf("过期 Token 应为 401, got %d", w.Code)
	}
}

func TestAdminJWTWrongSecret(t *testing.T) {
	signer := auth.NewJWTService("secret-a", time.Hour)
	token, err := signer.Generate("admin")
	if err != nil {
		t.Fatal(err)
	}

	verifier := auth.NewJWTService("secret-b", time.Hour)
	if w := doRequest(newTestRouter(verifier), "Bearer "+token); w.Code != http.StatusUnauthorized {
		t.Errorf("密钥不符应为 401, got %d", w.Code)
	}
}
