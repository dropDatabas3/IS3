package services

import (
	"testing"

	github_com_glebarez_sqlite "github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	categoriesClient "github.com/Guidotss/ucc-soft-arch-golang.git/src/clients/categories"
	"github.com/Guidotss/ucc-soft-arch-golang.git/src/domain/dtos/categories"
	"github.com/Guidotss/ucc-soft-arch-golang.git/src/model"
)

func setupCategoriesClientSQLite(t *testing.T) *categoriesClient.CategoriesClient {
	db, err := gorm.Open(github_com_glebarez_sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	// migrate required tables
	require.NoError(t, db.AutoMigrate(&model.Category{}))
	return categoriesClient.NewCategoryClient(db)
}

func TestCategoriesService_CreateAndFindAll(t *testing.T) {
	client := setupCategoriesClientSQLite(t)
	svc := NewCategoriesService(client)

	// create
	created, err := svc.CreateCategory(categories.CreateCategoryRequestDto{CategoryName: "Backend"})
	require.NoError(t, err)
	require.Equal(t, "Backend", created.CategoryName)
	require.NotEqual(t, created.CategoryId, model.Category{}.Id) // non-zero uuid

	// ensure persisted
	var inDB model.Category
	require.NoError(t, client.Db.First(&inDB, "id = ?", created.CategoryId).Error)
	require.Equal(t, "Backend", inDB.CategoryName)

	// find all mapping
	list, err := svc.FindAllCategories()
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, created.CategoryId, list[0].CategoryId)
	require.Equal(t, "Backend", list[0].CategoryName)
}
