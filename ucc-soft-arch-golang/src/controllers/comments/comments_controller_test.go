package comments

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	dto "github.com/Guidotss/ucc-soft-arch-golang.git/src/domain/dtos/comments"
	middlewares "github.com/Guidotss/ucc-soft-arch-golang.git/src/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type stubCommentsService struct {
	newResp dto.CommentRequestResponseDto
	newErr  error
	getResp dto.GetCommentsResponse
	getErr  error
	updResp dto.CommentRequestResponseDto
	updErr  error
}

func (s *stubCommentsService) NewComment(d dto.CommentRequestResponseDto) (dto.CommentRequestResponseDto, error) {
	return s.newResp, s.newErr
}
func (s *stubCommentsService) GetCourseComments(id uuid.UUID) (dto.GetCommentsResponse, error) {
	return s.getResp, s.getErr
}
func (s *stubCommentsService) UpdateComment(d dto.CommentRequestResponseDto) (dto.CommentRequestResponseDto, error) {
	return s.updResp, s.updErr
}

func TestCommentsController_NewComment_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubCommentsService{newResp: dto.CommentRequestResponseDto{Text: "ok"}}
	ctrl := NewCommentsController(svc)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.POST("/comments", ctrl.NewComment)
	req := httptest.NewRequest(http.MethodPost, "/comments", strings.NewReader(`{"text":"ok"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestCommentsController_GetCourseComments_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubCommentsService{}
	ctrl := NewCommentsController(svc)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.GET("/comments/:id", ctrl.GetCourseComments)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/comments/not-a-uuid", nil))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCommentsController_GetCourseComments_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubCommentsService{getResp: dto.GetCommentsResponse{{Text: "hi"}}}
	ctrl := NewCommentsController(svc)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.GET("/comments/:id", ctrl.GetCourseComments)
	id := uuid.New().String()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/comments/"+id, nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCommentsController_UpdateComment_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubCommentsService{}
	ctrl := NewCommentsController(svc)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.PUT("/comments", ctrl.UpdateComment)
	req := httptest.NewRequest(http.MethodPut, "/comments", strings.NewReader("{bad"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCommentsController_UpdateComment_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubCommentsService{updResp: dto.CommentRequestResponseDto{Text: "upd"}}
	ctrl := NewCommentsController(svc)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.PUT("/comments", ctrl.UpdateComment)
	req := httptest.NewRequest(http.MethodPut, "/comments", strings.NewReader(`{"text":"upd"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCommentsController_GetCourseComments_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubCommentsService{getErr: errors.New("boom")}
	ctrl := NewCommentsController(svc)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.GET("/comments/:id", ctrl.GetCourseComments)
	id := uuid.New().String()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/comments/"+id, nil))
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
