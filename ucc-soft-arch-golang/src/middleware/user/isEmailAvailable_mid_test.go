package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Guidotss/ucc-soft-arch-golang.git/src/domain/dtos/users"
	middlewares "github.com/Guidotss/ucc-soft-arch-golang.git/src/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// fake implementation of IUserService
type fakeUserService struct{ exists bool }

func (f fakeUserService) CreateUser(user users.RegisterRequest) (users.RegisterResponse, error) {
	return users.RegisterResponse{}, nil
}
func (f fakeUserService) GetUserById(id uuid.UUID) (users.GetUserDto, error) {
	return users.GetUserDto{}, nil
}
func (f fakeUserService) GetUserByEmail(email string) (users.GetUserDto, error) {
	if f.exists {
		return users.GetUserDto{Email: email}, nil
	}
	return users.GetUserDto{}, assertErr{}
}
func (f fakeUserService) UpdateUser(dto users.UpdateRequestDto) (users.UpdateResponseDto, error) {
	return users.UpdateResponseDto{}, nil
}

type assertErr struct{}

func (assertErr) Error() string { return "not found" }

func TestIsEmailAvailable_BlockedWhenExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.Use(func(c *gin.Context) { c.Set("Email", "taken@x.com"); c.Next() })
	r.Use(IsEmailAvailable(fakeUserService{exists: true}))
	r.GET("/path", func(c *gin.Context) { c.Status(http.StatusOK) })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/path", nil))
	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409 when email exists, got %d", w.Code)
	}
}

func TestIsEmailAvailable_PassesWhenNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("Email", "free@x.com"); c.Next() })
	r.Use(IsEmailAvailable(fakeUserService{exists: false}))
	r.GET("/path", func(c *gin.Context) { c.Status(http.StatusOK) })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/path", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected pass through, got %d", w.Code)
	}
}
