package courses

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"errors"

	domain "github.com/Guidotss/ucc-soft-arch-golang.git/src/domain/dtos/courses"
	middlewares "github.com/Guidotss/ucc-soft-arch-golang.git/src/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// stubCourseService is a minimal stub implementing services.ICourseService
type stubCourseService struct {
	findAllResp domain.GetAllCourses
	findAllErr  error
	findOneResp domain.GetCourseDto
	findOneErr  error
}

func (s *stubCourseService) CreateCourse(_ domain.CreateCoursesRequestDto) (domain.CreateCoursesResponseDto, error) {
	return domain.CreateCoursesResponseDto{}, nil
}
func (s *stubCourseService) FindAllCourses(_ string) (domain.GetAllCourses, error) {
	return s.findAllResp, s.findAllErr
}
func (s *stubCourseService) FindOneCourse(_ uuid.UUID) (domain.GetCourseDto, error) {
	if s.findOneErr != nil {
		return domain.GetCourseDto{}, s.findOneErr
	}
	return s.findOneResp, nil
}
func (s *stubCourseService) UpdateCourse(_ domain.UpdateRequestDto) (domain.UpdateResponseDto, error) {
	return domain.UpdateResponseDto{}, nil
}
func (s *stubCourseService) DeleteCourse(_ uuid.UUID) error { return nil }

func TestCourseController_GetAll_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubCourseService{findAllResp: domain.GetAllCourses{{CourseName: "Intro Go"}}}
	ctrl := NewCourseController(svc)

	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.GET("/courses", ctrl.GetAll)

	req := httptest.NewRequest(http.MethodGet, "/courses", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	// Validate body contains expected course name; avoid coupling to exact JSON field casing
	if !json.Valid(w.Body.Bytes()) {
		t.Fatalf("expected valid JSON body, got: %s", w.Body.String())
	}
	if !contains(w.Body.String(), "Intro Go") {
		t.Fatalf("expected response to include course name 'Intro Go', got: %s", w.Body.String())
	}
}

func TestCourseController_GetById_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubCourseService{}
	ctrl := NewCourseController(svc)

	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.GET("/courses/:id", ctrl.GetById)

	req := httptest.NewRequest(http.MethodGet, "/courses/not-a-uuid", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError && w.Code != http.StatusBadRequest { // gin Error middleware may map it
		t.Fatalf("expected 400/500 for invalid uuid, got %d", w.Code)
	}
}

func TestCourseController_GetById_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	expected := domain.GetCourseDto{CourseName: "Arch Soft"}
	svc := &stubCourseService{findOneResp: expected}
	ctrl := NewCourseController(svc)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.GET("/courses/:id", ctrl.GetById)
	id := uuid.New()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/courses/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !contains(w.Body.String(), "Arch Soft") {
		t.Fatalf("response missing course name: %s", w.Body.String())
	}
}

func TestCourseController_GetAll_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubCourseService{findAllErr: errors.New("boom")}
	ctrl := NewCourseController(svc)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.GET("/courses", ctrl.GetAll)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/courses", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 for service error, got %d", w.Code)
	}
}

// tiny helper to avoid importing another package just for substring
func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || (len(sub) > 0 && (func() bool { return stringIndex(s, sub) >= 0 })()))
}
func stringIndex(s, sub string) int {
	// naive scan; small strings in tests only
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
