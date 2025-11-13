package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	github_com_glebarez_sqlite "github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestAppRoutes_NoRouteReturnsError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// use a throwaway in-memory db for adapter wiring
	db, err := gorm.Open(github_com_glebarez_sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}

	AppRoutes(r, db)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)
	// default gin 404 since no error handler middleware is attached in AppRoutes
	require.Contains(t, w.Body.String(), "404 page not found")
}
