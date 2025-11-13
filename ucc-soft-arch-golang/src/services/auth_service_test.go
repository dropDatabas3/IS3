package services

import (
	"errors"
	"testing"

	userClient "github.com/Guidotss/ucc-soft-arch-golang.git/src/clients/users"
	userDtos "github.com/Guidotss/ucc-soft-arch-golang.git/src/domain/dtos/users"
	"github.com/Guidotss/ucc-soft-arch-golang.git/src/model"
	"github.com/Guidotss/ucc-soft-arch-golang.git/src/utils/bcrypt"
	"github.com/Guidotss/ucc-soft-arch-golang.git/src/utils/jwt"
	github_com_glebarez_sqlite "github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// fakeUserSvc implements IUserService for RefreshToken tests
type fakeUserSvc struct {
	user userDtos.GetUserDto
	err  error
}

func (f *fakeUserSvc) CreateUser(_ userDtos.RegisterRequest) (userDtos.RegisterResponse, error) {
	return userDtos.RegisterResponse{}, nil
}
func (f *fakeUserSvc) GetUserById(_ uuid.UUID) (userDtos.GetUserDto, error) { return f.user, f.err }
func (f *fakeUserSvc) GetUserByEmail(_ string) (userDtos.GetUserDto, error) { return f.user, f.err }
func (f *fakeUserSvc) UpdateUser(_ userDtos.UpdateRequestDto) (userDtos.UpdateResponseDto, error) {
	return userDtos.UpdateResponseDto{}, nil
}

func setupUsersClientWithSQLite(t *testing.T) *userClient.UsersClient {
	db, err := gorm.Open(github_com_glebarez_sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	// seed one user
	hashed, _ := bcrypt.HasPassword("secret")
	u := model.User{Id: uuid.New(), Email: "test@example.com", Password: hashed, Name: "Tester", Role: 1}
	if err := db.Create(&u).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return userClient.NewUsersClient(db)
}

func TestAuthService_Login_Success(t *testing.T) {
	client := setupUsersClientWithSQLite(t)
	// fake IUserService (not used by Login, but required by constructor)
	var us IUserService = &fakeUserSvc{}
	svc := NewAuthService(&us, client)

	dto := userDtos.LoginRequestDto{Email: "test@example.com", Password: "secret"}
	user, token, err := svc.Login(dto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "" {
		t.Fatalf("expected non-empty token")
	}
	if user.Email != "test@example.com" {
		t.Fatalf("unexpected user email: %s", user.Email)
	}
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	client := setupUsersClientWithSQLite(t)
	var us IUserService = &fakeUserSvc{}
	svc := NewAuthService(&us, client)

	dto := userDtos.LoginRequestDto{Email: "test@example.com", Password: "wrong"}
	_, _, err := svc.Login(dto)
	if err == nil {
		t.Fatalf("expected error for invalid credentials")
	}
}

func TestAuthService_RefreshToken_Success(t *testing.T) {
	// Prepare a user and token
	id := uuid.New()
	role := 1
	token := jwt.SignDocument(id, role)
	returned := userDtos.GetUserDto{Id: id, Email: "u@ex.com", Role: role, UserName: "U"}
	var us IUserService = &fakeUserSvc{user: returned, err: nil}
	// UsersClient not used by RefreshToken
	dummyClient := &userClient.UsersClient{}
	svc := NewAuthService(&us, dummyClient)

	user, newToken, err := svc.RefreshToken(token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if newToken == "" {
		t.Fatalf("expected new token")
	}
	if user.Id != id {
		t.Fatalf("expected same user id from user service")
	}
}

func TestAuthService_RefreshToken_Invalid(t *testing.T) {
	var us IUserService = &fakeUserSvc{}
	dummyClient := &userClient.UsersClient{}
	svc := NewAuthService(&us, dummyClient)
	_, _, err := svc.RefreshToken("invalid.token")
	if err == nil {
		t.Fatalf("expected error for invalid token")
	}
}

func TestAuthService_Login_EmailNotFound(t *testing.T) {
	// empty DB (no user seeded)
	db, err := gorm.Open(github_com_glebarez_sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	client := userClient.NewUsersClient(db)
	var us IUserService = &fakeUserSvc{}
	svc := NewAuthService(&us, client)
	_, _, err = svc.Login(userDtos.LoginRequestDto{Email: "missing@example.com", Password: "pw"})
	if err == nil {
		t.Fatalf("expected error when email not found")
	}
}

func TestAuthService_RefreshToken_UserLookupError(t *testing.T) {
	badErr := errors.New("db down")
	var us IUserService = &fakeUserSvc{err: badErr}
	svc := NewAuthService(&us, &userClient.UsersClient{})
	token := jwt.SignDocument(uuid.New(), 1)
	_, _, err := svc.RefreshToken(token)
	if err == nil {
		t.Fatalf("expected error when user service fails")
	}
}
