package services

import (
	"testing"

	github_com_glebarez_sqlite "github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	commentsClient "github.com/Guidotss/ucc-soft-arch-golang.git/src/clients/comments"
	dto "github.com/Guidotss/ucc-soft-arch-golang.git/src/domain/dtos/comments"
	"github.com/Guidotss/ucc-soft-arch-golang.git/src/model"
)

func setupCommentsClientSQLite(t *testing.T) *commentsClient.CommentsClient {
	db, err := gorm.Open(github_com_glebarez_sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	require.NoError(t, db.AutoMigrate(&model.User{}, &model.Course{}, &model.Comment{}))
	return commentsClient.NewCommentsClient(db)
}

func TestCommentsService_New_Get_Update(t *testing.T) {
	client := setupCommentsClientSQLite(t)
	svc := NewCommentsService(client)
	// seed dependencies
	u := model.User{Email: "c@e.com", Password: "p", Name: "Com"}
	cat := model.Category{CategoryName: "Cat"}
	c := model.Course{CourseName: "C1", CourseDescription: "d", CoursePrice: 1, CourseDuration: 1, CourseCapacity: 1, CourseInitDate: "2024", CategoryID: uuid.Nil}
	require.NoError(t, client.Db.Create(&u).Error)
	require.NoError(t, client.Db.Create(&cat).Error)
	c.CategoryID = cat.Id
	require.NoError(t, client.Db.Create(&c).Error)

	// new comment
	cr, err := svc.NewComment(dto.CommentRequestResponseDto{UserId: u.Id, CourseId: c.Id, Text: "hi"})
	require.NoError(t, err)
	require.Equal(t, "hi", cr.Text)

	// get comments
	list, err := svc.GetCourseComments(c.Id)
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, "hi", list[0].Text)

	// update comment
	ur, err := svc.UpdateComment(dto.CommentRequestResponseDto{UserId: u.Id, CourseId: c.Id, Text: "updated"})
	require.NoError(t, err)
	require.Equal(t, "updated", ur.Text)
}
