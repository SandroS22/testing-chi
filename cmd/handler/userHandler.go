package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/SandroS22/chi-test/cmd/request"
	"github.com/SandroS22/chi-test/cmd/response"
	"github.com/SandroS22/chi-test/cmd/util"
	"github.com/SandroS22/chi-test/db"
)

type UserHandler struct {
	DB *sql.DB
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req request.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := request.NewCreateUserRequest(req.Name, req.Email, req.Password, req.Role)

	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword
	dbUserID, err := db.InsertUser(h.DB, user.Name, user.Email, user.Password, user.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := response.NewCreateUserResponse(dbUserID, user.Name, user.Email, user.Password, user.Role)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := db.GetAllUsers(h.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
