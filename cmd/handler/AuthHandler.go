package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	customerrors "github.com/SandroS22/chi-test/cmd/CustomErrors"
	"github.com/SandroS22/chi-test/cmd/request"
	"github.com/SandroS22/chi-test/repository"
	"github.com/SandroS22/chi-test/service"
)

type AuthHandler struct {
	DB      *sql.DB
	Service service.AuthService
}

func (h *AuthHandler) getService() service.AuthService {
	if h.Service != nil {
		return h.Service
	}
	if h.DB != nil {
		repo := repository.NewUserRepository(h.DB)
		h.Service = service.NewAuthService(repo)
		return h.Service
	}
	return nil
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req request.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	svc := h.getService()
	if svc == nil {
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	res, err := svc.Login(req)
	if err != nil {
		if errors.Is(err, customerrors.ErrInvalidCredentials) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
