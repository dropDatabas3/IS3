package services

import (
	"testing"

	github_com_glebarez_sqlite "github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	ratingClient "github.com/Guidotss/ucc-soft-arch-golang.git/src/clients/rating"
	dto "github.com/Guidotss/ucc-soft-arch-golang.git/src/domain/dtos/rating"
	"github.com/Guidotss/ucc-soft-arch-golang.git/src/model"
)

func setupRatingClientSQLite(t *testing.T) *ratingClient.RatingClient {
	db, err := gorm.Open(github_com_glebarez_sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	require.NoError(t, db.AutoMigrate(&model.User{}, &model.Course{}, &model.Rating{}))
	return ratingClient.NewRatingClient(db)
}

func TestRatingService_New_Update_Get(t *testing.T) {
	client := setupRatingClientSQLite(t)
	svc := NewRatingService(client)
	// create deps
	u := model.User{Email: "u@e.com", Password: "p", Name: "U"}
	c := model.Course{CourseName: "Course", CourseDescription: "d", CoursePrice: 1, CourseDuration: 1, CourseInitDate: "2024", CourseCapacity: 1}
	require.NoError(t, client.Db.Create(&u).Error)
	require.NoError(t, client.Db.Create(&c).Error)

	// new rating
	r, err := svc.NewRating(dto.RatingRequestResponseDto{UserId: u.Id, CourseId: c.Id, Rating: 4})
	require.NoError(t, err)
	require.Equal(t, 4, r.Rating)

	// update rating
	updated, err := svc.UpdateRating(dto.RatingRequestResponseDto{UserId: u.Id, CourseId: c.Id, Rating: 5})
	require.NoError(t, err)
	require.Equal(t, 5, updated.Rating)

	// list ratings
	list, err := svc.GetRatings()
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, 5, list[0].Rating)
}
