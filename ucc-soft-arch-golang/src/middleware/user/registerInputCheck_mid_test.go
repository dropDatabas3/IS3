package user

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	middlewares "github.com/Guidotss/ucc-soft-arch-golang.git/src/middleware"
	"github.com/gin-gonic/gin"
)

func TestRegisterInputCheck_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.Use(RegisterInputCheckMiddleware())
	r.POST("/register", func(c *gin.Context) { c.Status(http.StatusOK) })
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestRegisterInputCheck_MissingFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.Use(RegisterInputCheckMiddleware())
	r.POST("/register", func(c *gin.Context) { c.Status(http.StatusOK) })
	body := bytes.NewBufferString(`{"email":"e@x.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestRegisterInputCheck_SetsContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RegisterInputCheckMiddleware())
	r.POST("/register", func(c *gin.Context) {
		if c.GetString("Email") == "" || c.GetString("Username") == "" || c.GetString("Password") == "" {
			t.Fatalf("expected context keys to be set")
		}
		c.Status(http.StatusOK)
	})
	body := bytes.NewBufferString(`{"email":"e@x.com","username":"john","password":"p"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
