package courses

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	domain "github.com/Guidotss/ucc-soft-arch-golang.git/src/domain/dtos/courses"
	middlewares "github.com/Guidotss/ucc-soft-arch-golang.git/src/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type fakeCourseService struct {
	updateResp domain.UpdateResponseDto
}

func (f *fakeCourseService) CreateCourse(_ domain.CreateCoursesRequestDto) (domain.CreateCoursesResponseDto, error) {
	return domain.CreateCoursesResponseDto{}, nil
}
func (f *fakeCourseService) FindAllCourses(_ string) (domain.GetAllCourses, error) { return nil, nil }
func (f *fakeCourseService) FindOneCourse(_ uuid.UUID) (domain.GetCourseDto, error) {
	return domain.GetCourseDto{}, nil
}
func (f *fakeCourseService) UpdateCourse(req domain.UpdateRequestDto) (domain.UpdateResponseDto, error) {
	// echo back some of the request for validation
	resp := f.updateResp
	resp.Id = req.Id
	if req.CourseName != nil {
		resp.CourseName = *req.CourseName
	}
	return resp, nil
}
func (f *fakeCourseService) DeleteCourse(_ uuid.UUID) error { return nil }

var _ interface {
	CreateCourse(domain.CreateCoursesRequestDto) (domain.CreateCoursesResponseDto, error)
	FindAllCourses(string) (domain.GetAllCourses, error)
	FindOneCourse(uuid.UUID) (domain.GetCourseDto, error)
	UpdateCourse(domain.UpdateRequestDto) (domain.UpdateResponseDto, error)
	DeleteCourse(uuid.UUID) error
} = (*fakeCourseService)(nil)

func TestCourseController_Create_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := NewCourseController(&fakeCourseService{})
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.POST("/courses/create", ctrl.Create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/courses/create", bytes.NewBufferString("not-json"))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 400/500, got %d", w.Code)
	}
}

func TestCourseController_Update_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := NewCourseController(&fakeCourseService{})
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.PUT("/courses/update/:id", ctrl.UpdateCourse)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/courses/update/"+uuid.New().String(), bytes.NewBufferString("not-json"))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 400/500, got %d", w.Code)
	}
}

func TestCourseController_Update_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &fakeCourseService{updateResp: domain.UpdateResponseDto{}}
	ctrl := NewCourseController(svc)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	// middleware to set courseId in context as string
	r.PUT("/courses/update/:id", func(c *gin.Context) { c.Set("courseId", uuid.New().String()) }, ctrl.UpdateCourse)
	body := map[string]any{"course_name": "Advanced Go"}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/courses/update/"+uuid.New().String(), bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (%s)", w.Code, w.Body.String())
	}
	if !json.Valid(w.Body.Bytes()) {
		t.Fatalf("expected json body")
	}
}

func TestCourseController_Delete_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := NewCourseController(&fakeCourseService{})
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.DELETE("/courses/:id", ctrl.DeleteCourse)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/courses/not-a-uuid", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 400/500, got %d", w.Code)
	}
}

func TestCourseController_Delete_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := NewCourseController(&fakeCourseService{})
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.DELETE("/courses/:id", ctrl.DeleteCourse)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/courses/"+uuid.New().String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
