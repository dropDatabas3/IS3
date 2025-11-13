package course

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCheckCourseId_MissingParam(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CheckCourseId())
	r.GET("/courses/", func(c *gin.Context) { c.Status(http.StatusOK) })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/courses/", nil))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCheckCourseId_SetsContextAndNext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CheckCourseId())
	r.GET("/courses/:id", func(c *gin.Context) {
		if _, exists := c.Get("courseId"); !exists {
			t.Fatalf("courseId not set")
		}
		c.Status(http.StatusOK)
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/courses/123", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
