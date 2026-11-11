package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SandroS22/chi-test/cmd/request"
	"github.com/SandroS22/chi-test/cmd/response"
	"github.com/SandroS22/chi-test/db"
	"github.com/SandroS22/chi-test/repository"
	"github.com/SandroS22/chi-test/service"
	"github.com/go-chi/chi/v5"
)

func setupTestHandler(t *testing.T) (*UserHandler, *sql.DB) {
	t.Helper()
	database, err := db.InitDB(":memory:")
	if err != nil {
		t.Fatalf("failed to init test db: %v", err)
	}

	repo := repository.NewUserRepository(database)
	svc := service.NewUserService(repo)
	handler := NewUserHandler(svc)
	return handler, database
}

func addChiURLParams(r *http.Request, params map[string]string) *http.Request {
	chiCtx := chi.NewRouteContext()
	for key, value := range params {
		chiCtx.URLParams.Add(key, value)
	}
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
}

func TestUserHandler_CreateUser(t *testing.T) {
	h, database := setupTestHandler(t)
	defer database.Close()

	body, _ := json.Marshal(request.CreateUserRequest{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "password123",
		Role:     "admin",
	})

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateUser(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", w.Code, w.Body.String())
	}

	var res response.CreateUserResponse
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if res.ID <= 0 || res.Name != "Alice" || res.Email != "alice@example.com" {
		t.Errorf("unexpected response: %+v", res)
	}
}

func TestUserHandler_CreateUser_BadRequest(t *testing.T) {
	h, database := setupTestHandler(t)
	defer database.Close()

	// Invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte("{invalid-json")))
	w := httptest.NewRecorder()
	h.CreateUser(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	// Missing required fields
	body, _ := json.Marshal(request.CreateUserRequest{Name: ""})
	req = httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	w = httptest.NewRecorder()
	h.CreateUser(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for empty fields, got %d", w.Code)
	}
}

func TestUserHandler_GetAllUsers(t *testing.T) {
	h, database := setupTestHandler(t)
	defer database.Close()

	// Create user first
	body, _ := json.Marshal(request.CreateUserRequest{
		Name:     "Bob",
		Email:    "bob@example.com",
		Password: "password",
		Role:     "user",
	})
	w := httptest.NewRecorder()
	h.CreateUser(w, httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body)))

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w = httptest.NewRecorder()
	h.GetAllUsers(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var users []response.UserResponse
	if err := json.Unmarshal(w.Body.Bytes(), &users); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
}

func TestUserHandler_GetUserByID(t *testing.T) {
	h, database := setupTestHandler(t)
	defer database.Close()

	// Create a user
	body, _ := json.Marshal(request.CreateUserRequest{
		Name:     "Charlie",
		Email:    "charlie@example.com",
		Password: "password",
		Role:     "user",
	})
	w := httptest.NewRecorder()
	h.CreateUser(w, httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body)))

	var created response.CreateUserResponse
	json.Unmarshal(w.Body.Bytes(), &created)

	// Valid ID
	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	req = addChiURLParams(req, map[string]string{"id": "1"})
	w = httptest.NewRecorder()
	h.GetUserByID(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var user response.UserResponse
	if err := json.Unmarshal(w.Body.Bytes(), &user); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if user.Name != "Charlie" || user.Email != "charlie@example.com" {
		t.Errorf("unexpected user data: %+v", user)
	}

	// Invalid ID
	req = httptest.NewRequest(http.MethodGet, "/users/abc", nil)
	req = addChiURLParams(req, map[string]string{"id": "abc"})
	w = httptest.NewRecorder()
	h.GetUserByID(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid id, got %d", w.Code)
	}

	// Not Found ID
	req = httptest.NewRequest(http.MethodGet, "/users/999", nil)
	req = addChiURLParams(req, map[string]string{"id": "999"})
	w = httptest.NewRecorder()
	h.GetUserByID(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404 for nonexistent id, got %d", w.Code)
	}
}

func TestUserHandler_UpdateUser(t *testing.T) {
	h, database := setupTestHandler(t)
	defer database.Close()

	// Create user
	body, _ := json.Marshal(request.CreateUserRequest{
		Name:     "David",
		Email:    "david@example.com",
		Password: "password",
		Role:     "user",
	})
	w := httptest.NewRecorder()
	h.CreateUser(w, httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body)))

	// Update user
	updateBody, _ := json.Marshal(request.UpdateUserRequest{
		Name:  "David New",
		Email: "david_new@example.com",
		Role:  "admin",
	})
	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader(updateBody))
	req = addChiURLParams(req, map[string]string{"id": "1"})
	w = httptest.NewRecorder()
	h.UpdateUser(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var updated response.UserResponse
	if err := json.Unmarshal(w.Body.Bytes(), &updated); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if updated.Name != "David New" || updated.Email != "david_new@example.com" || updated.Role != "admin" {
		t.Errorf("unexpected updated data: %+v", updated)
	}

	// Update with invalid ID
	req = httptest.NewRequest(http.MethodPut, "/users/abc", bytes.NewReader(updateBody))
	req = addChiURLParams(req, map[string]string{"id": "abc"})
	w = httptest.NewRecorder()
	h.UpdateUser(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	// Update nonexistent
	req = httptest.NewRequest(http.MethodPut, "/users/999", bytes.NewReader(updateBody))
	req = addChiURLParams(req, map[string]string{"id": "999"})
	w = httptest.NewRecorder()
	h.UpdateUser(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestUserHandler_DeleteUser(t *testing.T) {
	h, database := setupTestHandler(t)
	defer database.Close()

	// Create user
	body, _ := json.Marshal(request.CreateUserRequest{
		Name:     "Eva",
		Email:    "eva@example.com",
		Password: "password",
		Role:     "user",
	})
	w := httptest.NewRecorder()
	h.CreateUser(w, httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body)))

	// Delete user
	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	req = addChiURLParams(req, map[string]string{"id": "1"})
	w = httptest.NewRecorder()
	h.DeleteUser(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	// Delete again should return 404
	req = httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	req = addChiURLParams(req, map[string]string{"id": "1"})
	w = httptest.NewRecorder()
	h.DeleteUser(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404 on second delete, got %d", w.Code)
	}

	// Delete invalid ID
	req = httptest.NewRequest(http.MethodDelete, "/users/xyz", nil)
	req = addChiURLParams(req, map[string]string{"id": "xyz"})
	w = httptest.NewRecorder()
	h.DeleteUser(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 on invalid id, got %d", w.Code)
	}
}

func TestUserHandler_FallbackWithDB(t *testing.T) {
	database, err := db.InitDB(":memory:")
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	defer database.Close()

	// Legacy struct initialization with DB directly
	h := &UserHandler{DB: database}

	body, _ := json.Marshal(request.CreateUserRequest{
		Name:     "Fallback User",
		Email:    "fallback@example.com",
		Password: "password",
		Role:     "user",
	})
	w := httptest.NewRecorder()
	h.CreateUser(w, httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body)))

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}
