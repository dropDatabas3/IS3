package users

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	userDtos "github.com/Guidotss/ucc-soft-arch-golang.git/src/domain/dtos/users"
	middlewares "github.com/Guidotss/ucc-soft-arch-golang.git/src/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type stubUserService struct {
	createResp   userDtos.RegisterResponse
	createErr    error
	getEmailResp userDtos.GetUserDto
	getEmailErr  error
	updateResp   userDtos.UpdateResponseDto
	updateErr    error
}

func (s *stubUserService) CreateUser(r userDtos.RegisterRequest) (userDtos.RegisterResponse, error) {
	return s.createResp, s.createErr
}
func (s *stubUserService) GetUserById(id uuid.UUID) (userDtos.GetUserDto, error) {
	return userDtos.GetUserDto{}, nil
}
func (s *stubUserService) GetUserByEmail(email string) (userDtos.GetUserDto, error) {
	return s.getEmailResp, s.getEmailErr
}
func (s *stubUserService) UpdateUser(r userDtos.UpdateRequestDto) (userDtos.UpdateResponseDto, error) {
	return s.updateResp, s.updateErr
}

func TestUsersController_CreateUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubUserService{createResp: userDtos.RegisterResponse{Id: uuid.New(), Email: "new@ex.com"}}
	ctrl := NewUserController(svc)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.POST("/users", func(c *gin.Context) {
		c.Set("Email", "new@ex.com")
		c.Set("Username", "Neo")
		c.Set("Password", "pw")
		ctrl.CreateUser(c)
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/users", nil))
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestUsersController_FindByEmail_MissingEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := NewUserController(&stubUserService{})
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.GET("/users/email", ctrl.FindByEmail)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/users/email", nil))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUsersController_FindByEmail_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubUserService{getEmailResp: userDtos.GetUserDto{Email: "x@ex.com"}}
	ctrl := NewUserController(svc)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.GET("/users/email", func(c *gin.Context) { c.Set("email", "x@ex.com"); ctrl.FindByEmail(c) })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/users/email", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestUsersController_UpdateUser_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := NewUserController(&stubUserService{})
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.PUT("/users", func(c *gin.Context) { c.Set("userID", uuid.New()); ctrl.UpdateUser(c) })
	req := httptest.NewRequest(http.MethodPut, "/users", strings.NewReader("{bad"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUsersController_UpdateUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubUserService{updateResp: userDtos.UpdateResponseDto{Email: "u@ex.com"}}
	ctrl := NewUserController(svc)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.PUT("/users", func(c *gin.Context) { c.Set("userID", uuid.New()); ctrl.UpdateUser(c) })
	req := httptest.NewRequest(http.MethodPut, "/users", strings.NewReader(`{"email":"u@ex.com"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestUsersController_FindByEmail_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubUserService{getEmailErr: errors.New("fail")}
	ctrl := NewUserController(svc)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.GET("/users/email", func(c *gin.Context) { c.Set("email", "x@ex.com"); ctrl.FindByEmail(c) })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/users/email", nil))
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
