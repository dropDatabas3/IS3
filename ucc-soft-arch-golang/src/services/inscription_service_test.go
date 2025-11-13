package services

import (
	"testing"

	github_com_glebarez_sqlite "github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	inscClient "github.com/Guidotss/ucc-soft-arch-golang.git/src/clients/inscriptos"
	dto "github.com/Guidotss/ucc-soft-arch-golang.git/src/domain/dtos/inscription"
	"github.com/Guidotss/ucc-soft-arch-golang.git/src/model"
)

func setupInscriptosClientSQLite(t *testing.T) *inscClient.InscriptosClient {
	db, err := gorm.Open(github_com_glebarez_sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	require.NoError(t, db.AutoMigrate(&model.User{}, &model.Category{}, &model.Course{}, &model.Inscripto{}))
	return inscClient.NewInscriptionClient(db)
}

func seedUser(t *testing.T, db *gorm.DB, email, name string) model.User {
	u := model.User{Email: email, Password: "x", Name: name}
	require.NoError(t, db.Create(&u).Error)
	return u
}

func TestInscriptionService_Enroll_And_Queries(t *testing.T) {
	client := setupInscriptosClientSQLite(t)
	svc := NewInscriptionService(client)

	cat := model.Category{CategoryName: "Cloud"}
	require.NoError(t, client.Db.Create(&cat).Error)
	course := model.Course{CourseName: "Azure 101", CourseDescription: "intro", CoursePrice: 1, CourseDuration: 1, CourseCapacity: 10, CourseInitDate: "2024-01-01", CourseState: true, CategoryID: cat.Id}
	require.NoError(t, client.Db.Create(&course).Error)
	user := seedUser(t, client.Db, "a@b.com", "Alice")

	// Enroll
	eresp, err := svc.Enroll(dto.EnrollRequestResponseDto{CourseId: course.Id, UserId: user.Id})
	require.NoError(t, err)
	require.Equal(t, course.Id, eresp.CourseId)
	require.Equal(t, user.Id, eresp.UserId)

	// GetMyStudents
	students, err := svc.GetMyStudents(course.Id)
	require.NoError(t, err)
	require.Len(t, students, 1)
	require.Equal(t, "Alice", students[0].UserName)

	// IsUserEnrolled
	enrolled, err := svc.IsUserEnrolled(user.Id, course.Id)
	require.NoError(t, err)
	require.True(t, enrolled)

	// CourseExist
	exist, err := svc.CourseExist(course.Id)
	require.NoError(t, err)
	require.True(t, exist)
	noexist, err := svc.CourseExist(uuid.New())
	require.NoError(t, err)
	require.False(t, noexist)
}
