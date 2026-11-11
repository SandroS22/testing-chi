package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/SandroS22/chi-test/cmd/request"
	"github.com/SandroS22/chi-test/cmd/response"
	"github.com/SandroS22/chi-test/repository"
	"github.com/SandroS22/chi-test/service"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	DB      *sql.DB
	Service service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		Service: userService,
	}
}

func (h *UserHandler) getService() service.UserService {
	if h.Service != nil {
		return h.Service
	}
	if h.DB != nil {
		repo := repository.NewUserRepository(h.DB)
		h.Service = service.NewUserService(repo)
		return h.Service
	}
	return nil
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req request.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	svc := h.getService()
	if svc == nil {
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	res, err := svc.CreateUser(req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserData) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	svc := h.getService()
	if svc == nil {
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	users, err := svc.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	svc := h.getService()
	if svc == nil {
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	user, err := svc.GetUserByID(id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	h.GetUserByID(w, r)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	var req request.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	svc := h.getService()
	if svc == nil {
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	user, err := svc.UpdateUser(id, req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrInvalidUserData) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	svc := h.getService()
	if svc == nil {
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	if err := svc.DeleteUser(id); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.NewResponse(http.StatusOK, "User deleted successfully", nil))
}
