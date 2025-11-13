package inscriptos

import (
	"testing"

	github_com_glebarez_sqlite "github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/Guidotss/ucc-soft-arch-golang.git/src/model"
)

func setupInscriptosDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(github_com_glebarez_sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	require.NoError(t, db.AutoMigrate(&model.User{}, &model.Course{}, &model.Inscripto{}, &model.Category{}))
	return db
}

func TestInscriptosClient_Enroll_IsUserEnrolled_CourseExist_GetMyStudents(t *testing.T) {
	db := setupInscriptosDB(t)
	c := NewInscriptionClient(db)

	// seed category, course, user
	cat := model.Category{CategoryName: "Backend"}
	require.NoError(t, db.Create(&cat).Error)
	course := model.Course{CourseName: "Golang", CourseDescription: "intro", CoursePrice: 10, CourseDuration: 8, CourseInitDate: "2025-01-01", CourseState: false, CourseCapacity: 30, CourseImage: "img", CategoryID: cat.Id}
	require.NoError(t, db.Create(&course).Error)
	u := model.User{Name: "Alice", Avatar: "pic.png", Email: "a@b.com", Password: "x"}
	require.NoError(t, db.Create(&u).Error)

	// CourseExist should be true after creating course
	ok, err := c.CourseExist(course.Id)
	require.NoError(t, err)
	require.True(t, ok)

	// Enroll
	ins := model.Inscripto{UserId: u.Id, CourseId: course.Id}
	_, err = c.Enroll(ins)
	require.NoError(t, err)

	// IsUserEnrolled
	enrolled, err := c.IsUserEnrolled(u.Id, course.Id)
	require.NoError(t, err)
	require.True(t, enrolled)

	// GetMyStudents
	students, err := c.GetMyStudents(course.Id)
	require.NoError(t, err)
	require.Len(t, students, 1)
	require.Equal(t, "Alice", students[0].Name)
	require.Equal(t, "pic.png", students[0].Avatar)
	// Note: raw scan key alias casing can differ; avoid strict assertion on Id value here
}

func TestInscriptosClient_IsUserEnrolled_And_CourseExist_False(t *testing.T) {
	db := setupInscriptosDB(t)
	c := NewInscriptionClient(db)

	// Create a course but no enrollment
	cat := model.Category{CategoryName: "Backend"}
	require.NoError(t, db.Create(&cat).Error)
	course := model.Course{CourseName: "Golang", CourseDescription: "intro", CoursePrice: 10, CourseDuration: 8, CourseInitDate: "2025-01-01", CourseState: false, CourseCapacity: 30, CourseImage: "img", CategoryID: cat.Id}
	require.NoError(t, db.Create(&course).Error)
	u := model.User{Name: "Bob", Avatar: "a.png", Email: "b@b.com", Password: "x"}
	require.NoError(t, db.Create(&u).Error)

	enrolled, err := c.IsUserEnrolled(u.Id, course.Id)
	require.NoError(t, err)
	require.False(t, enrolled)

	exists, err := c.CourseExist(uuid.New())
	require.NoError(t, err)
	require.False(t, exists)
}

func TestInscriptosClient_GetMyCourses(t *testing.T) {
	db := setupInscriptosDB(t)
	c := NewInscriptionClient(db)

	// seed category, course, user, and enrollment
	cat := model.Category{CategoryName: "Backend"}
	require.NoError(t, db.Create(&cat).Error)
	course := model.Course{CourseName: "Golang", CourseDescription: "intro", CoursePrice: 15.5, CourseDuration: 12, CourseInitDate: "2025-01-01", CourseState: true, CourseCapacity: 100, CourseImage: "img", CategoryID: cat.Id}
	require.NoError(t, db.Create(&course).Error)
	u := model.User{Name: "Carol", Avatar: "c.png", Email: "c@c.com", Password: "x"}
	require.NoError(t, db.Create(&u).Error)
	ins := model.Inscripto{UserId: u.Id, CourseId: course.Id}
	require.NoError(t, db.Create(&ins).Error)

	courses, err := c.GetMyCourses(u.Id)
	require.NoError(t, err)
	require.Len(t, courses, 1)
	got := courses[0]
	require.Equal(t, "Golang", got.CourseName)
	require.Equal(t, "Backend", got.Category.CategoryName)
	// Type-normalized fields
	require.InDelta(t, 15.5, got.CoursePrice, 0.0001)
	require.Equal(t, 12, got.CourseDuration)
	require.True(t, got.CourseState)
	require.Equal(t, 100, got.CourseCapacity)
}

func TestInscriptosClient_GetMyStudents_NotFound(t *testing.T) {
	db := setupInscriptosDB(t)
	c := NewInscriptionClient(db)

	// Seed only a course, but no enrollments
	cat := model.Category{CategoryName: "X"}
	require.NoError(t, db.Create(&cat).Error)
	course := model.Course{CourseName: "NoStudents", CourseDescription: "", CoursePrice: 0, CourseDuration: 1, CourseInitDate: "2025-01-01", CourseState: false, CourseCapacity: 10, CourseImage: "", CategoryID: cat.Id}
	require.NoError(t, db.Create(&course).Error)

	_, err := c.GetMyStudents(course.Id)
	require.Error(t, err)
}

func TestInscriptosClient_Enroll_DBError(t *testing.T) {
	// Fresh DB without migrations to force Create error
	rawDB, err := gorm.Open(github_com_glebarez_sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	c := NewInscriptionClient(rawDB)

	_, err = c.Enroll(model.Inscripto{})
	require.Error(t, err)
}

func TestInscriptosClient_GetMyCourses_DBError(t *testing.T) {
	rawDB, err := gorm.Open(github_com_glebarez_sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	c := NewInscriptionClient(rawDB)

	_, err = c.GetMyCourses(uuid.New())
	require.Error(t, err)
}

func TestInscriptosClient_GetMyStudents_DBError(t *testing.T) {
	rawDB, err := gorm.Open(github_com_glebarez_sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	c := NewInscriptionClient(rawDB)

	_, err = c.GetMyStudents(uuid.New())
	require.Error(t, err)
}

func TestInscriptosClient_IsUserEnrolled_DBError(t *testing.T) {
	// No migrations -> Count on non-existent table should error
	rawDB, err := gorm.Open(github_com_glebarez_sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	c := NewInscriptionClient(rawDB)

	enrolled, err := c.IsUserEnrolled(uuid.New(), uuid.New())
	require.Error(t, err)
	require.False(t, enrolled)
}

func TestInscriptosClient_CourseExist_DBError(t *testing.T) {
	// No migrations -> Count on non-existent table should error
	rawDB, err := gorm.Open(github_com_glebarez_sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	c := NewInscriptionClient(rawDB)

	exists, err := c.CourseExist(uuid.New())
	require.Error(t, err)
	require.False(t, exists)
}
