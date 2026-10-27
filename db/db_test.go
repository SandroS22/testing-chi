package db

import (
	"database/sql"
	"testing"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	// Use in-memory SQLite for testing
	database, err := InitDB(":memory:")
	if err != nil {
		t.Fatalf("failed to init test db: %v", err)
	}
	return database
}

func TestInitDB(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	if err := database.Ping(); err != nil {
		t.Fatalf("failed to ping db: %v", err)
	}
}

func TestInsertAndGetUser(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	id, err := InsertUser(database, "John Doe", "johndoe", "admin")
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	if id <= 0 {
		t.Fatalf("expected valid id, got %d", id)
	}

	user, err := GetUserByID(database, id)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}

	if user.Name != "John Doe" || user.User != "johndoe" || user.Role != "admin" {
		t.Errorf("unexpected user data: %+v", user)
	}
}

func TestGetAllUsers(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	_, err := InsertUser(database, "User One", "user1", "member")
	if err != nil {
		t.Fatalf("failed to insert user 1: %v", err)
	}

	_, err = InsertUser(database, "User Two", "user2", "admin")
	if err != nil {
		t.Fatalf("failed to insert user 2: %v", err)
	}

	users, err := GetAllUsers(database)
	if err != nil {
		t.Fatalf("failed to get all users: %v", err)
	}

	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}

	if users[0].User != "user1" || users[1].User != "user2" {
		t.Errorf("unexpected users order/data: %+v", users)
	}
}

func TestUpdateUser(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	id, err := InsertUser(database, "Original Name", "orig", "viewer")
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	err = UpdateUser(database, User{
		ID:   id,
		Name: "Updated Name",
		User: "updated",
		Role: "editor",
	})
	if err != nil {
		t.Fatalf("failed to update user: %v", err)
	}

	updated, err := GetUserByID(database, id)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}

	if updated.Name != "Updated Name" || updated.User != "updated" || updated.Role != "editor" {
		t.Errorf("user was not updated correctly: %+v", updated)
	}
}

func TestDeleteUser(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	id, err := InsertUser(database, "To Delete", "delete_me", "guest")
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	if err := DeleteUser(database, id); err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}

	user, err := GetUserByID(database, id)
	if err == nil {
		t.Fatalf("expected error when fetching deleted user, got: %+v", user)
	}
}
