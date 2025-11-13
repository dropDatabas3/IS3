package categories

import (
	"testing"

	github_com_glebarez_sqlite "github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/Guidotss/ucc-soft-arch-golang.git/src/model"
)

func setupCatDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(github_com_glebarez_sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	require.NoError(t, db.AutoMigrate(&model.Category{}))
	return db
}

func TestCategoriesClient_CreateAndGetAll(t *testing.T) {
	db := setupCatDB(t)
	c := NewCategoryClient(db)

	cat := model.Category{CategoryName: "Backend"}
	created, err := c.Create(cat)
	require.NoError(t, err)
	require.Equal(t, "Backend", created.CategoryName)

	list, err := c.GetAll()
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, "Backend", list[0].CategoryName)
}

func TestCategoriesClient_Create_ErrorUnexpected(t *testing.T) {
	db := setupCatDB(t)
	// Force write into a non-existent table to trigger default error mapping
	c := NewCategoryClient(db.Table("does_not_exist"))

	_, err := c.Create(model.Category{CategoryName: "X"})
	require.Error(t, err)
}

func TestCategoriesClient_GetAll_DBError(t *testing.T) {
	db := setupCatDB(t)
	// Querying a non-existent table should surface a DB error branch
	c := NewCategoryClient(db.Table("does_not_exist"))

	_, err := c.GetAll()
	require.Error(t, err)
}
