package main

import (
	"log"
	"net/http"
	"userServer/internal/handlers/http/handlerOperator"

	"userServer/internal/handlers/http/handlersPlan"
	"userServer/internal/handlers/http/handlersUser"
	"userServer/internal/repository/postgres/repositoryOperetor"
	"userServer/internal/repository/postgres/repositoryPlan"
	"userServer/internal/repository/postgres/repositoryUser"
	"userServer/internal/service/serverOperetor"
	"userServer/internal/service/servicePlan"
	"userServer/internal/service/serviceUser"

	"github.com/gorilla/mux"
)

func userInit() *handlersUser.UserHandler {
	return handlersUser.NewUserHandler(serviceUser.NewUserService(repositoryUser.NewUserRepository()))
}

func planInit() *handlersPlan.PlanHandler {
	return handlersPlan.NewPlanHandler(servicePlan.NewPlanService(repositoryPlan.NewPlanRepository()))
}

func operatorInit() *handlerOperator.OperatorHandler {
	return handlerOperator.NewOperatorHandler(serverOperetor.NewOperatorService(repositoryOperetor.NewOperetorRepository()))
}

func main() {
	router := mux.NewRouter()
	userHandler := userInit()
	planHandler := planInit()
	operatorHandler := operatorInit()

	router.HandleFunc("/register-user", userHandler.Create).Methods("POST") // TODO
	router.HandleFunc("/update-user", userHandler.Update).Methods("PUT")    // TODO
	router.HandleFunc("/get-user/{id}", userHandler.Get).Methods("GET")     // TODO

	router.HandleFunc("/get-plan/{id}", planHandler.Get).Methods("GET") // TODO
	router.HandleFunc("/get-plans", planHandler.GetAll).Methods("GET")  // TODO

	router.HandleFunc("/get-operator/{id}", operatorHandler.Get).Methods("GET") // TODO
	router.HandleFunc("/get-operators", operatorHandler.GetAll).Methods("GET")  // TODO

	log.Println("Server starting on :8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
