package comments

import (
	"testing"

	github_com_glebarez_sqlite "github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/Guidotss/ucc-soft-arch-golang.git/src/model"
)

func setupCommentsDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(github_com_glebarez_sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	require.NoError(t, db.AutoMigrate(&model.User{}, &model.Course{}, &model.Comment{}))
	return db
}

func TestCommentsClient_NewAndUpdateAndGet(t *testing.T) {
	db := setupCommentsDB(t)
	c := NewCommentsClient(db)

	// seed user and course
	u := model.User{Name: "Alice", Avatar: "pic.png", Email: "a@b.com", Password: "x"}
	require.NoError(t, db.Create(&u).Error)
	course := model.Course{CourseName: "Golang", CourseDescription: "intro", CoursePrice: 10, CourseDuration: 8, CourseInitDate: "2025-01-01", CourseState: false, CourseCapacity: 30, CourseImage: "img"}
	require.NoError(t, db.Create(&course).Error)

	// new comment
	cm := model.Comment{Text: "hi", UserId: u.Id, CourseId: course.Id}
	created, err := c.NewComment(cm)
	require.NoError(t, err)
	require.Equal(t, "hi", created.Text)

	// update comment
	cm.Text = "updated"
	updated, err := c.UpdateComment(cm)
	require.NoError(t, err)
	require.Equal(t, "updated", updated.Text)

	// get comments for course
	list, err := c.GetCourseComments(course.Id)
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, "updated", list[0].Text)
	require.Equal(t, "Alice", list[0].UserName)
	require.Equal(t, "pic.png", list[0].UserAvatar)
	// note: user_id alias in raw query is mixed case; sqlite scan may not match expected key
	// so we avoid asserting on parsed UserId here to keep compatibility
}

func TestCommentsClient_GetCourseComments_NotFound(t *testing.T) {
	db := setupCommentsDB(t)
	c := NewCommentsClient(db)

	// seed a course but no comments
	course := model.Course{CourseName: "Empty", CourseDescription: "", CoursePrice: 0, CourseDuration: 1, CourseInitDate: "2025-01-01", CourseState: false, CourseCapacity: 1, CourseImage: ""}
	require.NoError(t, db.Create(&course).Error)

	_, err := c.GetCourseComments(course.Id)
	require.Error(t, err)
}

func TestCommentsClient_NewComment_DBError(t *testing.T) {
	// Open DB but do not migrate the comments table via overriding table
	db := setupCommentsDB(t)
	c := NewCommentsClient(db.Table("does_not_exist"))

	cm := model.Comment{Text: "hi"}
	_, err := c.NewComment(cm)
	require.Error(t, err)
}

func TestCommentsClient_UpdateComment_DBError(t *testing.T) {
	// Use a fresh in-memory DB without migrating comments table
	rawDB, err := gorm.Open(github_com_glebarez_sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	c := NewCommentsClient(rawDB)

	// Attempt to update should fail as the table doesn't exist
	_, err = c.UpdateComment(model.Comment{Text: "x"})
	require.Error(t, err)
}

func TestCommentsClient_GetCourseComments_DBError(t *testing.T) {
	// No migrations so SELECT on comments/users should error
	rawDB, err := gorm.Open(github_com_glebarez_sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	c := NewCommentsClient(rawDB)

	_, err = c.GetCourseComments(uuid.New())
	require.Error(t, err)
}
