package main

import (
	"log"
	"net/http"
	handlers "userServer/internal/handlers/http"
	"userServer/internal/repository/postgres"
	"userServer/internal/service"
)

func main() {

	repo := postgres.NewUserRepository()
	userService := service.NewUserService(repo)
	userHandler := handlers.NewUserHandler(userService)

	http.HandleFunc("/user", userHandler.CreateUser)   //TODO
	http.HandleFunc("/user/", userHandler.GetUser)     //TODO
	http.HandleFunc("/userUp", userHandler.UpdateUser) //TODO

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
