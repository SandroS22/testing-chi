package handler

import (
	"encoding/json"
	"net/http"

	"github.com/SandroS22/chi-test/cmd/response"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	testResponse := response.NewTestResponse(1, "John Doe", "johndoe", "admin")
	response := response.NewResponse(http.StatusOK, "Test Handler", testResponse)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
