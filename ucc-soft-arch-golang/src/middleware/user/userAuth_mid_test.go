package user

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	middlewares "github.com/Guidotss/ucc-soft-arch-golang.git/src/middleware"
	"github.com/Guidotss/ucc-soft-arch-golang.git/src/utils/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.Use(AuthMiddleware())
	r.GET("/secure", func(c *gin.Context) { c.Status(http.StatusOK) })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/secure", nil))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.Use(AuthMiddleware())
	r.GET("/secure", func(c *gin.Context) { c.Status(http.StatusOK) })
	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	req.Header.Set("Authorization", "Token xxx")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestAuthMiddleware_SetsUserIDOnSuccess(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.Use(AuthMiddleware())
	r.GET("/secure", func(c *gin.Context) {
		if _, exists := c.Get("userID"); !exists {
			t.Fatalf("userID not set")
		}
		c.Status(http.StatusOK)
	})
	uid := uuid.New()
	token := jwt.SignDocument(uid, 0)
	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
