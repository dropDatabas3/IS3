package rating

import (
	"testing"

	github_com_glebarez_sqlite "github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/Guidotss/ucc-soft-arch-golang.git/src/model"
)

func setupRatingDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(github_com_glebarez_sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	require.NoError(t, db.AutoMigrate(&model.User{}, &model.Course{}, &model.Rating{}))
	return db
}

func TestRatingClient_CRUD(t *testing.T) {
	db := setupRatingDB(t)
	c := NewRatingClient(db)

	// seed user and course
	u := model.User{Name: "Alice", Email: "a@b.com", Password: "x"}
	require.NoError(t, db.Create(&u).Error)
	course := model.Course{CourseName: "Golang", CourseDescription: "intro", CoursePrice: 10, CourseDuration: 8, CourseInitDate: "2025-01-01", CourseState: false, CourseCapacity: 30, CourseImage: "img"}
	require.NoError(t, db.Create(&course).Error)

	r := model.Rating{UserId: u.Id, CourseId: course.Id, Rating: 4}
	created, err := c.NewRating(r)
	require.NoError(t, err)
	require.Equal(t, 4, created.Rating)

	r.Rating = 5
	updated, err := c.UpdateRating(r)
	require.NoError(t, err)
	require.Equal(t, 5, updated.Rating)

	list, err := c.GetRatings()
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, uuid.Nil != list[0].UserId, true)
}
