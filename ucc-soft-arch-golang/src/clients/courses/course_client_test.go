package courses

import (
	"testing"

	github_com_glebarez_sqlite "github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	customError "github.com/Guidotss/ucc-soft-arch-golang.git/src/domain/errors"
	"github.com/Guidotss/ucc-soft-arch-golang.git/src/model"
)

func setupCoursesDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(github_com_glebarez_sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	require.NoError(t, db.AutoMigrate(&model.Course{}, &model.Category{}, &model.User{}, &model.Rating{}))
	return db
}

func TestCourseClient_Create_Update_Delete(t *testing.T) {
	db := setupCoursesDB(t)
	c := NewCourseClient(db)

	// seed category
	cat := model.Category{CategoryName: "Backend"}
	require.NoError(t, db.Create(&cat).Error)

	course := model.Course{CourseName: "Golang", CourseDescription: "intro", CoursePrice: 10, CourseDuration: 8, CourseInitDate: "2025-01-01", CourseState: false, CourseCapacity: 30, CourseImage: "img", CategoryID: cat.Id}
	created, err := c.Create(course)
	require.NoError(t, err)
	require.Equal(t, "Golang", created.CourseName)

	// update some fields
	created.CourseName = "Golang 2"
	updated, err := c.UpdateCourse(created)
	require.NoError(t, err)
	require.Equal(t, "Golang 2", updated.CourseName)

	// delete
	err = c.DeleteCourse(created.Id)
	require.NoError(t, err)
}

func TestCourseClient_GetAll_And_GetById(t *testing.T) {
	db := setupCoursesDB(t)
	c := NewCourseClient(db)

	// seed data
	cat := model.Category{CategoryName: "Backend"}
	require.NoError(t, db.Create(&cat).Error)
	usr := model.User{Email: "a@b.com", Password: "pw", Name: "A"}
	require.NoError(t, db.Create(&usr).Error)

	course := model.Course{CourseName: "Golang", CourseDescription: "intro", CoursePrice: 20, CourseDuration: 8, CourseInitDate: "2025-01-01", CourseState: true, CourseCapacity: 25, CourseImage: "img", CategoryID: cat.Id}
	require.NoError(t, db.Create(&course).Error)

	// one rating so AVG = rating
	r := model.Rating{CourseId: course.Id, UserId: usr.Id, Rating: 4}
	require.NoError(t, db.Create(&r).Error)

	// GetAll without filter
	all, err := c.GetAll("")
	require.NoError(t, err)
	require.Len(t, all, 1)
	require.Equal(t, "Golang", all[0].CourseName)
	require.Equal(t, "Backend", all[0].Category.CategoryName)
	require.InDelta(t, 4.0, all[0].RatingAvg, 0.0001)
	require.Equal(t, true, all[0].CourseState)

	// GetById
	got, err := c.GetById(course.Id)
	require.NoError(t, err)
	require.Equal(t, course.Id, got.Id)
	require.Equal(t, "Golang", got.CourseName)
	require.Equal(t, true, got.CourseState)
	require.InDelta(t, 4.0, got.RatingAvg, 0.0001)
}

func TestCourseClient_FilterAndNotFoundBranches(t *testing.T) {
	db := setupCoursesDB(t)
	c := NewCourseClient(db)

	// seed category and a course without ratings
	cat := model.Category{CategoryName: "Backend"}
	require.NoError(t, db.Create(&cat).Error)
	course := model.Course{CourseName: "Rust", CourseDescription: "systems", CoursePrice: 30, CourseDuration: 10, CourseInitDate: "2025-02-01", CourseState: true, CourseCapacity: 50, CourseImage: "img", CategoryID: cat.Id}
	require.NoError(t, db.Create(&course).Error)

	// 1) GetAll with a filter triggers the filtered query branch.
	//    The client uses LOWER(...) LIKE LOWER(?) so it should work in sqlite.
	list, err := c.GetAll("Rust")
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, "Rust", list[0].CourseName)

	// 2) GetById for a course without ratings should now return NOT_FOUND (no join hit)
	_, err = c.GetById(course.Id)
	require.Error(t, err)
}

func TestCourseClient_Create_DuplicateName_ErrorMapping(t *testing.T) {
	db := setupCoursesDB(t)
	c := NewCourseClient(db)

	// seed category
	cat := model.Category{CategoryName: "Backend"}
	require.NoError(t, db.Create(&cat).Error)

	// create one course
	course := model.Course{CourseName: "UniqueName", CourseDescription: "d", CoursePrice: 10, CourseDuration: 5, CourseInitDate: "2025-01-01", CourseState: true, CourseCapacity: 10, CourseImage: "img", CategoryID: cat.Id}
	_, err := c.Create(course)
	require.NoError(t, err)

	// attempt to create another with the same unique name triggers sqlite UNIQUE constraint error
	_, err = c.Create(course)
	require.Error(t, err)

	// our client maps specific postgres duplicate error text; sqlite will hit default branch (UNEXPECTED_ERROR)
	ce, ok := err.(*customError.Error)
	require.True(t, ok)
	require.Equal(t, "UNEXPECTED_ERROR", ce.Code)
}

func TestCourseClient_DBError_Paths(t *testing.T) {
	// Use a raw in-memory DB with no migrations to force SQL errors across methods
	rawDB, err := gorm.Open(github_com_glebarez_sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	c := NewCourseClient(rawDB)

	// Create should fail (no courses table) -> mapped to UNEXPECTED_ERROR for sqlite
	_, err = c.Create(model.Course{})
	require.Error(t, err)

	// Update should fail (no courses table)
	_, err = c.UpdateCourse(model.Course{})
	require.Error(t, err)

	// Delete should fail (no courses table)
	err = c.DeleteCourse(model.Course{}.Id)
	require.Error(t, err)

	// GetAll should fail (no tables) -> DB_ERROR mapping
	_, err = c.GetAll("")
	require.Error(t, err)

	// GetById should fail (no tables) -> DB_ERROR mapping
	_, err = c.GetById(model.Course{}.Id)
	require.Error(t, err)
}
