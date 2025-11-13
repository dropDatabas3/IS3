package inscriptions

import (
	"net/http"
	"net/http/httptest"
	"testing"

	courseDto "github.com/Guidotss/ucc-soft-arch-golang.git/src/domain/dtos/courses"
	inDto "github.com/Guidotss/ucc-soft-arch-golang.git/src/domain/dtos/inscription"
	middlewares "github.com/Guidotss/ucc-soft-arch-golang.git/src/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type stubInscriptionService struct {
	enrollResp     inDto.EnrollRequestResponseDto
	enrollErr      error
	myCourses      courseDto.GetAllCourses
	myCoursesErr   error
	students       inDto.StudentsInCourse
	studentsErr    error
	isEnrolled     bool
	isEnrolledErr  error
	courseExist    bool
	courseExistErr error
}

func (s *stubInscriptionService) Enroll(d inDto.EnrollRequestResponseDto) (inDto.EnrollRequestResponseDto, error) {
	return s.enrollResp, s.enrollErr
}
func (s *stubInscriptionService) GetMyCourses(id uuid.UUID) (courseDto.GetAllCourses, error) {
	return s.myCourses, s.myCoursesErr
}
func (s *stubInscriptionService) GetMyStudents(id uuid.UUID) (inDto.StudentsInCourse, error) {
	return s.students, s.studentsErr
}
func (s *stubInscriptionService) IsUserEnrolled(u, c uuid.UUID) (bool, error) {
	return s.isEnrolled, s.isEnrolledErr
}
func (s *stubInscriptionService) CourseExist(id uuid.UUID) (bool, error) {
	return s.courseExist, s.courseExistErr
}

func TestInscriptionController_Create_MissingIDs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := NewInscriptionController(&stubInscriptionService{})
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.POST("/inscriptions", ctrl.Create)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/inscriptions", nil))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestInscriptionController_Create_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubInscriptionService{enrollResp: inDto.EnrollRequestResponseDto{}}
	ctrl := NewInscriptionController(svc)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.POST("/inscriptions", ctrl.Create)
	req := httptest.NewRequest(http.MethodPost, "/inscriptions", nil)
	uid := uuid.New().String()
	cid := uuid.New().String()
	req = req.WithContext(req.Context())
	// set headers into context via gin.Set? we need to use c.Set in middleware: emulate by wrapping handler
	r = gin.New()
	r.Use(middlewares.ErrorHandler())
	r.POST("/inscriptions", func(c *gin.Context) { c.Set("userID", uid); c.Set("courseID", cid); ctrl.Create(c) })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestInscriptionController_GetMyCourses_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubInscriptionService{myCourses: courseDto.GetAllCourses{}}
	ctrl := NewInscriptionController(svc)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.GET("/mine", func(c *gin.Context) { c.Set("userID", uuid.New()); ctrl.GetMyCourses(c) })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/mine", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestInscriptionController_GetMyStudents_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := NewInscriptionController(&stubInscriptionService{})
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.GET("/students/:cid", ctrl.GetMyStudents)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/students/not-uuid", nil))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestInscriptionController_GetMyStudents_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubInscriptionService{students: inDto.StudentsInCourse{}}
	ctrl := NewInscriptionController(svc)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.GET("/students/:cid", ctrl.GetMyStudents)
	id := uuid.New().String()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/students/"+id, nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestInscriptionController_IsAlreadyEnrolled_True(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubInscriptionService{isEnrolled: true}
	ctrl := NewInscriptionController(svc)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.GET("/enrolled/:cid", func(c *gin.Context) {
		c.Set("userID", uuid.New().String())
		c.Params = append(c.Params, gin.Param{Key: "cid", Value: uuid.New().String()})
		ctrl.IsAlredyEnrolled(c)
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/enrolled/"+uuid.New().String(), nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 (controller continues after error), got %d", w.Code)
	}
}

func TestInscriptionController_IsAlreadyEnrolled_False(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubInscriptionService{isEnrolled: false}
	ctrl := NewInscriptionController(svc)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.GET("/enrolled/:cid", func(c *gin.Context) {
		c.Set("userID", uuid.New().String())
		c.Params = append(c.Params, gin.Param{Key: "cid", Value: uuid.New().String()})
		ctrl.IsAlredyEnrolled(c)
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/enrolled/"+uuid.New().String(), nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 when not enrolled, got %d", w.Code)
	}
}
