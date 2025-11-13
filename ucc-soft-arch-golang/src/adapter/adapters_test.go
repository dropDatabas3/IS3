package adapter

import (
	"testing"

	github_com_glebarez_sqlite "github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/Guidotss/ucc-soft-arch-golang.git/src/model"
)

func setupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(github_com_glebarez_sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	// migrate commonly used tables to ensure relationships exist if exercised
	require.NoError(t, db.AutoMigrate(&model.User{}, &model.Category{}, &model.Course{}, &model.Comment{}, &model.Rating{}, &model.Inscripto{}))
	return db
}

func TestUserAdapter(t *testing.T) {
	db := setupDB(t)
	ctrl, svc := UserAdapter(db)
	require.NotNil(t, ctrl)
	require.NotNil(t, svc)
}

func TestCourseAdapter(t *testing.T) {
	db := setupDB(t)
	ctrl := CourseAdapter(db)
	require.NotNil(t, ctrl)
}

func TestCategoryAdapter(t *testing.T) {
	db := setupDB(t)
	ctrl := CategoryAdapter(db)
	require.NotNil(t, ctrl)
}

func TestCommentAdapter(t *testing.T) {
	db := setupDB(t)
	ctrl := CommentAdapter(db)
	require.NotNil(t, ctrl)
}

func TestInscriptionsAdapter(t *testing.T) {
	db := setupDB(t)
	ctrl, svc := InscriptionsAdapter(db)
	require.NotNil(t, ctrl)
	require.NotNil(t, svc)
}

func TestAuthAdapter(t *testing.T) {
	db := setupDB(t)
	ctrl := AuthAdapter(db)
	require.NotNil(t, ctrl)
}

func TestRatingAdapter(t *testing.T) {
	db := setupDB(t)
	ctrl := RatingAdapter(db)
	require.NotNil(t, ctrl)
}
