package main

import (
	"log"
	"net/http"
	handlers "userServer/internal/handlers/http"
	"userServer/internal/repository/postgres"
	"userServer/internal/service"
)

func main() {

	repo := postgres.NewMemoryRepo()
	userService := service.NewUserService(repo)
	userHandler := handlers.NewUserHandler(userService)

	http.HandleFunc("/createUser", userHandler.CreateUser)
	http.HandleFunc("/users/", userHandler.GetUser)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
