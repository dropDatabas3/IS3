package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	userDtos "github.com/Guidotss/ucc-soft-arch-golang.git/src/domain/dtos/users"
	middlewares "github.com/Guidotss/ucc-soft-arch-golang.git/src/middleware"
	"github.com/Guidotss/ucc-soft-arch-golang.git/src/services"
	"github.com/gin-gonic/gin"
)

// stubAuthService implements services.IAuthService behavior we need
type stubAuthService struct {
	loginUser    userDtos.GetUserDto
	loginToken   string
	loginErr     error
	refreshUser  userDtos.GetUserDto
	refreshToken string
	refreshErr   error
}

func (s *stubAuthService) Login(dto userDtos.LoginRequestDto) (userDtos.GetUserDto, string, error) {
	return s.loginUser, s.loginToken, s.loginErr
}
func (s *stubAuthService) RefreshToken(token string) (userDtos.GetUserDto, string, error) {
	return s.refreshUser, s.refreshToken, s.refreshErr
}

func makeAuthController(s *stubAuthService) *AuthController {
	var as services.IAuthService = s
	return NewAuthController(&as)
}

func TestAuthController_Login_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubAuthService{loginUser: userDtos.GetUserDto{Email: "x@ex.com"}, loginToken: "tok"}
	ctrl := makeAuthController(svc)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.POST("/login", ctrl.Login)
	body := `{"email":"x@ex.com","password":"pw"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "tok") {
		t.Fatalf("expected token in body: %s", w.Body.String())
	}
}

func TestAuthController_Login_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubAuthService{}
	ctrl := makeAuthController(svc)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.POST("/login", ctrl.Login)
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("{bad json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestAuthController_RefreshToken_MissingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubAuthService{}
	ctrl := makeAuthController(svc)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.GET("/refresh", ctrl.RefreshToken)
	req := httptest.NewRequest(http.MethodGet, "/refresh", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestAuthController_RefreshToken_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubAuthService{refreshUser: userDtos.GetUserDto{Email: "r@ex.com"}, refreshToken: "newtok"}
	ctrl := makeAuthController(svc)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.GET("/refresh", ctrl.RefreshToken)
	req := httptest.NewRequest(http.MethodGet, "/refresh", nil)
	req.Header.Set("Authorization", "Bearer old")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var js map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &js); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if js["token"] != "newtok" {
		t.Fatalf("expected newtok token, got %v", js["token"])
	}
}

func TestAuthController_RefreshToken_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubAuthService{refreshErr: errors.New("fail")}
	ctrl := makeAuthController(svc)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.GET("/refresh", ctrl.RefreshToken)
	req := httptest.NewRequest(http.MethodGet, "/refresh", nil)
	req.Header.Set("Authorization", "Bearer old")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError && w.Code != http.StatusBadRequest {
		t.Fatalf("expected error status, got %d", w.Code)
	}
}
