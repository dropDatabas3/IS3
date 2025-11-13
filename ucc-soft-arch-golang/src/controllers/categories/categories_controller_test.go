package categories

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	catDtos "github.com/Guidotss/ucc-soft-arch-golang.git/src/domain/dtos/categories"
	middlewares "github.com/Guidotss/ucc-soft-arch-golang.git/src/middleware"
	"github.com/gin-gonic/gin"
)

type stubCategoriesService struct {
	createResp catDtos.CreateCategoryResponseDto
	createErr  error
	allResp    catDtos.GetAllCategories
	allErr     error
}

func (s *stubCategoriesService) CreateCategory(d catDtos.CreateCategoryRequestDto) (catDtos.CreateCategoryResponseDto, error) {
	return s.createResp, s.createErr
}
func (s *stubCategoriesService) FindAllCategories() (catDtos.GetAllCategories, error) {
	return s.allResp, s.allErr
}

func TestCategoriesController_Create_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubCategoriesService{createResp: catDtos.CreateCategoryResponseDto{CategoryName: "Dev"}}
	ctrl := NewCategoriesController(svc)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.POST("/categories", ctrl.Create)
	req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(`{"category_name":"Dev"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestCategoriesController_Create_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubCategoriesService{}
	ctrl := NewCategoriesController(svc)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.POST("/categories", ctrl.Create)
	req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader("{bad json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCategoriesController_GetAll_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubCategoriesService{allResp: catDtos.GetAllCategories{{CategoryName: "A"}}}
	ctrl := NewCategoriesController(svc)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.GET("/categories", ctrl.GetAll)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/categories", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCategoriesController_GetAll_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubCategoriesService{allErr: errors.New("fail")}
	ctrl := NewCategoriesController(svc)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.GET("/categories", ctrl.GetAll)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/categories", nil))
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
