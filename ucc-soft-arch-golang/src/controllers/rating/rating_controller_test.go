package rating

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	dto "github.com/Guidotss/ucc-soft-arch-golang.git/src/domain/dtos/rating"
	middlewares "github.com/Guidotss/ucc-soft-arch-golang.git/src/middleware"
	"github.com/gin-gonic/gin"
)

type stubRatingService struct {
	newResp dto.RatingRequestResponseDto
	newErr  error
	updResp dto.RatingRequestResponseDto
	updErr  error
	list    dto.RatingsResponse
	listErr error
}

func (s *stubRatingService) NewRating(d dto.RatingRequestResponseDto) (dto.RatingRequestResponseDto, error) {
	return s.newResp, s.newErr
}
func (s *stubRatingService) UpdateRating(d dto.RatingRequestResponseDto) (dto.RatingRequestResponseDto, error) {
	return s.updResp, s.updErr
}
func (s *stubRatingService) GetRatings() (dto.RatingsResponse, error) { return s.list, s.listErr }

func TestRatingController_NewRating_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := NewRatingController(&stubRatingService{newResp: dto.RatingRequestResponseDto{Rating: 5}})
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.POST("/ratings", ctrl.NewRating)
	req := httptest.NewRequest(http.MethodPost, "/ratings", strings.NewReader(`{"rating":5}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestRatingController_NewRating_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := NewRatingController(&stubRatingService{})
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.POST("/ratings", ctrl.NewRating)
	req := httptest.NewRequest(http.MethodPost, "/ratings", strings.NewReader("{bad"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestRatingController_UpdateRating_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := NewRatingController(&stubRatingService{updResp: dto.RatingRequestResponseDto{Rating: 3}})
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.PUT("/ratings", ctrl.UpdateRating)
	req := httptest.NewRequest(http.MethodPut, "/ratings", strings.NewReader(`{"rating":3}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestRatingController_UpdateRating_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := NewRatingController(&stubRatingService{})
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.PUT("/ratings", ctrl.UpdateRating)
	req := httptest.NewRequest(http.MethodPut, "/ratings", strings.NewReader("{bad"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestRatingController_GetRatings_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := NewRatingController(&stubRatingService{list: dto.RatingsResponse{{Rating: 1}}})
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.GET("/ratings", ctrl.GetRatings)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ratings", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestRatingController_GetRatings_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := NewRatingController(&stubRatingService{listErr: errors.New("fail")})
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.GET("/ratings", ctrl.GetRatings)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ratings", nil))
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
