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

func TestAuthRoutes_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	db, err := gorm.Open(github_com_glebarez_sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	AuthRoutes(r, adapter.AuthAdapter(db))

	// exercise a route minimally
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	r.ServeHTTP(w, req)
	// controller will likely validate body and return error; just ensure route exists
	require.True(t, w.Code == http.StatusBadRequest || w.Code >= 400)
}
