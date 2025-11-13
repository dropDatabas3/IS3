package enroll

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	courseDto "github.com/Guidotss/ucc-soft-arch-golang.git/src/domain/dtos/courses"
	dto "github.com/Guidotss/ucc-soft-arch-golang.git/src/domain/dtos/inscription"
	"github.com/Guidotss/ucc-soft-arch-golang.git/src/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// fakeInscriptionService implements services.IInscriptionService for middleware tests
type fakeInscriptionService struct {
	courseExists bool
	isEnrolled   bool
}

func (f *fakeInscriptionService) Enroll(d dto.EnrollRequestResponseDto) (dto.EnrollRequestResponseDto, error) {
	return dto.EnrollRequestResponseDto{}, nil
}
func (f *fakeInscriptionService) GetMyCourses(id uuid.UUID) (courseDto.GetAllCourses, error) {
	return nil, nil
}
func (f *fakeInscriptionService) GetMyStudents(id uuid.UUID) (dto.StudentsInCourse, error) {
	return nil, nil
}
func (f *fakeInscriptionService) IsUserEnrolled(userID uuid.UUID, courseID uuid.UUID) (bool, error) {
	return f.isEnrolled, nil
}
func (f *fakeInscriptionService) CourseExist(course_id uuid.UUID) (bool, error) {
	return f.courseExists, nil
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestCourseExistMiddleware_InvalidJSON(t *testing.T) {
	r := setupRouter()
	svc := &fakeInscriptionService{courseExists: false}
	r.POST("/enroll", CourseExist(svc))

	req := httptest.NewRequest(http.MethodPost, "/enroll", bytes.NewBufferString("not-json"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCourseExistMiddleware_InvalidUUID(t *testing.T) {
	r := setupRouter()
	svc := &fakeInscriptionService{courseExists: false}
	r.POST("/enroll", CourseExist(svc))

	payload := map[string]string{"course_id": "not-a-uuid"}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/enroll", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCourseExistMiddleware_NotExists(t *testing.T) {
	r := setupRouter()
	svc := &fakeInscriptionService{courseExists: false}
	r.POST("/enroll", CourseExist(svc))

	payload := map[string]string{"course_id": uuid.New().String()}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/enroll", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCourseExistMiddleware_SetsContextAndNext(t *testing.T) {
	r := setupRouter()
	svc := &fakeInscriptionService{courseExists: true}
	r.POST("/enroll", CourseExist(svc), func(c *gin.Context) {
		// verify context and return 204
		v, ok := c.Get("courseID")
		if !ok || v == "" {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusNoContent)
	})

	payload := map[string]string{"course_id": uuid.New().String()}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/enroll", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)
}

func TestIsAlreadyEnrollMiddleware_MissingUserID(t *testing.T) {
	r := setupRouter()
	svc := &fakeInscriptionService{isEnrolled: false}
	r.POST("/enroll", IsAlredyEnroll(svc))

	req := httptest.NewRequest(http.MethodPost, "/enroll", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestIsAlreadyEnrollMiddleware_MissingCourseID(t *testing.T) {
	r := setupRouter()
	svc := &fakeInscriptionService{isEnrolled: false}
	r.POST("/enroll", func(c *gin.Context) {
		c.Set("userID", uuid.New())
	}, IsAlredyEnroll(svc))

	req := httptest.NewRequest(http.MethodPost, "/enroll", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestIsAlreadyEnrollMiddleware_InvalidCourseIDFormat(t *testing.T) {
	r := setupRouter()
	svc := &fakeInscriptionService{isEnrolled: false}
	r.POST("/enroll", func(c *gin.Context) {
		c.Set("userID", uuid.New())
		c.Set("courseID", "not-uuid")
	}, IsAlredyEnroll(svc))

	req := httptest.NewRequest(http.MethodPost, "/enroll", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestIsAlreadyEnrollMiddleware_AlreadyEnrolled(t *testing.T) {
	r := setupRouter()
	svc := &fakeInscriptionService{isEnrolled: true}
	r.POST("/enroll", func(c *gin.Context) {
		c.Set("userID", uuid.New())
		c.Set("courseID", uuid.New().String())
	}, IsAlredyEnroll(svc))

	req := httptest.NewRequest(http.MethodPost, "/enroll", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestIsAlreadyEnrollMiddleware_PassesWhenNotEnrolled(t *testing.T) {
	r := setupRouter()
	svc := &fakeInscriptionService{isEnrolled: false}
	r.POST("/enroll", func(c *gin.Context) {
		c.Set("userID", uuid.New())
		c.Set("courseID", uuid.New().String())
	}, IsAlredyEnroll(svc), func(c *gin.Context) { c.Status(http.StatusNoContent) })

	req := httptest.NewRequest(http.MethodPost, "/enroll", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)
}

// Ensure fakeInscriptionService satisfies the interface at compile time
var _ services.IInscriptionService = (*fakeInscriptionService)(nil)
