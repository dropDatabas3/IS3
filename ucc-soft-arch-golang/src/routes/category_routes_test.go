package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	github_com_glebarez_sqlite "github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/Guidotss/ucc-soft-arch-golang.git/src/adapter"
)

func TestCategoryRoutes_RegisterAndGetAll(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	db, err := gorm.Open(github_com_glebarez_sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	CategoriesRoutes(r, adapter.CategoryAdapter(db))

	// call GET /categories should be 200 with empty slice
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
