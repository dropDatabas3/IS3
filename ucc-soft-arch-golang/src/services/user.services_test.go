package services

import (
	"testing"

	github_com_glebarez_sqlite "github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	usersClient "github.com/Guidotss/ucc-soft-arch-golang.git/src/clients/users"
	userDto "github.com/Guidotss/ucc-soft-arch-golang.git/src/domain/dtos/users"
	"github.com/Guidotss/ucc-soft-arch-golang.git/src/model"
	"github.com/Guidotss/ucc-soft-arch-golang.git/src/utils/bcrypt"
)

func setupUsersClientSQLite(t *testing.T) *usersClient.UsersClient {
	db, err := gorm.Open(github_com_glebarez_sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	require.NoError(t, db.AutoMigrate(&model.User{}))
	return usersClient.NewUsersClient(db)
}

func TestUserService_Create_Get_Update(t *testing.T) {
	client := setupUsersClientSQLite(t)
	svcInterface := NewUserService(client)
	svc := svcInterface.(*UserService)

	// Create
	reg, err := svc.CreateUser(userDto.RegisterRequest{
		Email:    "a@b.com",
		Password: "secret",
		Username: "alice",
		Avatar:   "pic.png",
	})
	require.NoError(t, err)
	require.Equal(t, "a@b.com", reg.Email)
	require.Equal(t, "alice", reg.Username)
	require.Equal(t, "pic.png", reg.Avatar)
	require.NotEqual(t, uuid.Nil, reg.Id)

	// Ensure password hashed in DB
	var inDB model.User
	require.NoError(t, client.Db.First(&inDB, "id = ?", reg.Id).Error)
	require.NotEqual(t, "secret", inDB.Password)
	require.True(t, bcrypt.ComparePassword("secret", inDB.Password))

	// Get by id
	byID, err := svc.GetUserById(reg.Id)
	require.NoError(t, err)
	require.Equal(t, reg.Id, byID.Id)
	require.Equal(t, "a@b.com", byID.Email)
	require.Equal(t, "alice", byID.UserName)

	// Get by email
	byEmail, err := svc.GetUserByEmail("a@b.com")
	require.NoError(t, err)
	require.Equal(t, byID, byEmail)

	// Update (username + avatar)
	upd, err := svc.UpdateUser(userDto.UpdateRequestDto{
		Id:       reg.Id,
		Username: "alice2",
		Avatar:   "pic2.png",
	})
	require.NoError(t, err)
	require.Equal(t, "alice2", upd.Username)
	require.Equal(t, "pic2.png", upd.Avatar)

	// Update password path (ensure hashed)
	_, err = svc.UpdateUser(userDto.UpdateRequestDto{
		Id:       reg.Id,
		Password: "newsecret",
	})
	require.NoError(t, err)
	require.NoError(t, client.Db.First(&inDB, "id = ?", reg.Id).Error)
	require.True(t, bcrypt.ComparePassword("newsecret", inDB.Password))
}
