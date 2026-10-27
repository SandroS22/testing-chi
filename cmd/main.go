package main

import (
	"log"

	"github.com/SandroS22/chi-test/cmd/config"
	"github.com/SandroS22/chi-test/cmd/handler"
	"github.com/SandroS22/chi-test/router"
	"github.com/go-chi/chi/v5"
)

func main() {
	config := config.LoadConfig("config.yml")
	r := router.NewRouter()
	registerRoutes(r)
	log.Println("Starting server on :" + config.Port)
	router.StartServer(r, ":"+config.Port)
}

func registerRoutes(r *chi.Mux) {
	router.AddRoute(r, "GET", "/hello", handler.HelloHandler)
	router.AddRoute(r, "GET", "/test", handler.TestHandler)
}
