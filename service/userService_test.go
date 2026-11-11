package service

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/SandroS22/chi-test/cmd/request"
	"github.com/SandroS22/chi-test/cmd/util"
	"github.com/SandroS22/chi-test/db"
	"github.com/SandroS22/chi-test/repository"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	database, err := db.InitDB(":memory:")
	if err != nil {
		t.Fatalf("failed to init test db: %v", err)
	}
	return database
}

func TestUserService_CreateUser(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	repo := repository.NewUserRepository(database)
	svc := NewUserService(repo)

	// Valid creation
	req := request.CreateUserRequest{
		Name:     "Bob",
		Email:    "bob@example.com",
		Password: "plainPassword123",
		Role:     "admin",
	}

	res, err := svc.CreateUser(req)
	if err != nil {
		t.Fatalf("unexpected error creating user: %v", err)
	}
	if res.ID <= 0 {
		t.Fatalf("expected valid id, got %d", res.ID)
	}
	if res.Name != "Bob" || res.Email != "bob@example.com" || res.Role != "admin" {
		t.Errorf("unexpected response: %+v", res)
	}
	if !util.CheckPasswordHash("plainPassword123", res.Password) {
		t.Errorf("password was not hashed properly")
	}

	// Missing required fields
	invalidReq := request.CreateUserRequest{
		Name: "",
	}
	_, err = svc.CreateUser(invalidReq)
	if !errors.Is(err, ErrInvalidUserData) {
		t.Errorf("expected ErrInvalidUserData, got %v", err)
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	repo := repository.NewUserRepository(database)
	svc := NewUserService(repo)

	created, err := svc.CreateUser(request.CreateUserRequest{
		Name:     "Charlie",
		Email:    "charlie@example.com",
		Password: "password",
		Role:     "user",
	})
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	user, err := svc.GetUserByID(created.ID)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}
	if user.Name != "Charlie" || user.Email != "charlie@example.com" {
		t.Errorf("unexpected user data: %+v", user)
	}

	// Not found
	_, err = svc.GetUserByID(9999)
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserService_GetAllUsers(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	repo := repository.NewUserRepository(database)
	svc := NewUserService(repo)

	// Initially empty
	users, err := svc.GetAllUsers()
	if err != nil {
		t.Fatalf("failed to get users: %v", err)
	}
	if len(users) != 0 {
		t.Errorf("expected 0 users, got %d", len(users))
	}

	_, _ = svc.CreateUser(request.CreateUserRequest{
		Name:     "User1",
		Email:    "u1@example.com",
		Password: "pass",
		Role:     "role1",
	})
	_, _ = svc.CreateUser(request.CreateUserRequest{
		Name:     "User2",
		Email:    "u2@example.com",
		Password: "pass",
		Role:     "role2",
	})

	users, err = svc.GetAllUsers()
	if err != nil {
		t.Fatalf("failed to get users: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	repo := repository.NewUserRepository(database)
	svc := NewUserService(repo)

	created, err := svc.CreateUser(request.CreateUserRequest{
		Name:     "Old Name",
		Email:    "old@example.com",
		Password: "oldpassword",
		Role:     "user",
	})
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	updated, err := svc.UpdateUser(created.ID, request.UpdateUserRequest{
		Name:     "New Name",
		Email:    "new@example.com",
		Role:     "admin",
		Password: "newpassword123",
	})
	if err != nil {
		t.Fatalf("failed to update user: %v", err)
	}

	if updated.Name != "New Name" || updated.Email != "new@example.com" || updated.Role != "admin" {
		t.Errorf("unexpected updated data: %+v", updated)
	}

	// Update non-existent user
	_, err = svc.UpdateUser(9999, request.UpdateUserRequest{Name: "Nobody"})
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	repo := repository.NewUserRepository(database)
	svc := NewUserService(repo)

	created, err := svc.CreateUser(request.CreateUserRequest{
		Name:     "Delete Me",
		Email:    "del@example.com",
		Password: "password",
		Role:     "guest",
	})
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	if err := svc.DeleteUser(created.ID); err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}

	// Verify not found
	_, err = svc.GetUserByID(created.ID)
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}

	// Delete non-existent user
	err = svc.DeleteUser(created.ID)
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound on second delete, got %v", err)
	}
}
