package repository

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/SandroS22/chi-test/db"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	database, err := db.InitDB(":memory:")
	if err != nil {
		t.Fatalf("failed to init test db: %v", err)
	}
	return database
}

func TestUserRepository_CRUD(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	repo := NewUserRepository(database)

	// 1. Create
	user := &db.User{
		Name:     "Alice",
		User:     "alice@example.com",
		Password: "hashedpassword123",
		Role:     "admin",
	}

	id, err := repo.Create(user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	if id <= 0 {
		t.Fatalf("expected valid id, got %d", id)
	}

	// 2. GetByID
	found, err := repo.GetByID(id)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}
	if found.Name != "Alice" || found.User != "alice@example.com" || found.Role != "admin" {
		t.Errorf("unexpected user data: %+v", found)
	}

	// 3. GetAll
	users, err := repo.GetAll()
	if err != nil {
		t.Fatalf("failed to get all users: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}

	// 4. Update
	updatedUser := &db.User{
		ID:       id,
		Name:     "Alice Updated",
		User:     "alice_updated@example.com",
		Password: "newhashedpassword",
		Role:     "user",
	}
	if err := repo.Update(updatedUser); err != nil {
		t.Fatalf("failed to update user: %v", err)
	}

	foundUpdated, err := repo.GetByID(id)
	if err != nil {
		t.Fatalf("failed to get updated user: %v", err)
	}
	if foundUpdated.Name != "Alice Updated" || foundUpdated.User != "alice_updated@example.com" || foundUpdated.Role != "user" {
		t.Errorf("unexpected updated user data: %+v", foundUpdated)
	}

	// 5. Delete
	if err := repo.Delete(id); err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}

	// Verify not found after delete
	_, err = repo.GetByID(id)
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserRepository_NotFoundCases(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	repo := NewUserRepository(database)

	// GetByID not found
	_, err := repo.GetByID(999)
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound on GetByID, got %v", err)
	}

	// Update not found
	err = repo.Update(&db.User{ID: 999, Name: "Nonexistent"})
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound on Update, got %v", err)
	}

	// Delete not found
	err = repo.Delete(999)
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound on Delete, got %v", err)
	}
}
