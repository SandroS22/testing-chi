package main

import (
	"database/sql"
	"log"

	"github.com/SandroS22/chi-test/cmd/config"
	"github.com/SandroS22/chi-test/cmd/handler"
	"github.com/SandroS22/chi-test/db"
	"github.com/SandroS22/chi-test/router"
	"github.com/go-chi/chi/v5"
)

func main() {
	config := config.LoadConfig("config.yml")
	dba, err := db.InitDB("application.db")
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	db.CreateTables(dba)
	defer dba.Close()
	r := router.NewRouter()
	registerRoutes(r)
	registerUserRoutes(r, dba)
	log.Println("Starting server on :" + config.Port)
	router.StartServer(r, ":"+config.Port)
}

func registerRoutes(r *chi.Mux) {
	router.AddRoute(r, "GET", "/hello", handler.HelloHandler)
	router.AddRoute(r, "GET", "/test", handler.TestHandler)
}

func registerUserRoutes(r *chi.Mux, db *sql.DB) {
	userHandler := createHandler(db)
	router.AddRoute(r, "POST", "/users", userHandler.CreateUser)
	router.AddRoute(r, "GET", "/users", userHandler.GetAllUsers)
}

func createHandler(db *sql.DB) *handler.UserHandler {
	return &handler.UserHandler{DB: db}
}
