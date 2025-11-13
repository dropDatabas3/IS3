package admin

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt"
)

// minimal custom claims to build token (matching utils/jwt if needed)
type testClaims struct {
	jwt.StandardClaims
	Role int
}

func makeToken(t *testing.T, role int, secret []byte) string {
	claims := testClaims{Role: role}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString(secret)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return s
}

func TestAdminAuthMiddleware_MissingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AdminAuthMiddleware())
	r.GET("/secure", func(c *gin.Context) { c.Status(http.StatusOK) })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/secure", nil))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAdminAuthMiddleware_InvalidFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AdminAuthMiddleware())
	r.GET("/secure", func(c *gin.Context) { c.Status(http.StatusOK) })
	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	req.Header.Set("Authorization", "BadToken")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAdminAuthMiddleware_ForbiddenRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := []byte("") // secret expected from env loader; using empty matches test scenario
	token := makeToken(t, 0, secret)
	r := gin.New()
	r.Use(AdminAuthMiddleware())
	r.GET("/secure", func(c *gin.Context) { c.Status(http.StatusOK) })
	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for role 0, got %d", w.Code)
	}
}

func TestAdminAuthMiddleware_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := []byte("")
	token := makeToken(t, 1, secret)
	r := gin.New()
	r.Use(AdminAuthMiddleware())
	r.GET("/secure", func(c *gin.Context) { c.Status(http.StatusOK) })
	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 success, got %d", w.Code)
	}
}

// ensure we don't accidentally allow whitespace issues
func TestAdminAuthMiddleware_TrimBearer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := []byte("")
	token := makeToken(t, 1, secret)
	r := gin.New()
	r.Use(AdminAuthMiddleware())
	r.GET("/secure", func(c *gin.Context) { c.Status(http.StatusOK) })
	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	req.Header.Set("Authorization", "Bearer  "+token) // two spaces -> should fail format split
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for malformed bearer spacing, got %d", w.Code)
	}
	// Now send a properly spaced header
	req2 := httptest.NewRequest(http.MethodGet, "/secure", nil)
	req2.Header.Set("Authorization", strings.TrimSpace("Bearer "+token))
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200 with proper spacing, got %d", w2.Code)
	}
}
