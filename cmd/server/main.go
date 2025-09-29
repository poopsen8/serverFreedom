package main

import (
	"log"
	"net/http"
	"userServer/internal/application/background"
	"userServer/internal/handlers/http/handlerOperator"
	"userServer/internal/handlers/http/handlerSubscription"

	"userServer/internal/handlers/http/handlersPlan"
	"userServer/internal/handlers/http/handlersUser"
	"userServer/internal/repository/postgres/repositoryOperetor"
	"userServer/internal/repository/postgres/repositoryPlan"
	"userServer/internal/repository/postgres/repositorySubscription"
	"userServer/internal/repository/postgres/repositoryUser"

	"userServer/internal/service/serviceOperetor"
	"userServer/internal/service/servicePlan"
	"userServer/internal/service/serviceSubscription"
	"userServer/internal/service/serviceUser"

	"github.com/gorilla/mux"
)

func initLayers() (*handlersUser.UserHandler, *handlersPlan.PlanHandler, *handlerOperator.OperatorHandler, *handlerSubscription.SubscriptionHandler, *background.TaskService) {
	operatorHandler := handlerOperator.NewOperatorHandler(serviceOperetor.NewOperatorService(repositoryOperetor.NewOperetorRepository()))
	userHandler := handlersUser.NewUserHandler(serviceUser.NewUserService(repositoryUser.NewUserRepository(), *serviceOperetor.NewOperatorService(repositoryOperetor.NewOperetorRepository())))
	planHandler := handlersPlan.NewPlanHandler(servicePlan.NewPlanService(repositoryPlan.NewPlanRepository()))
	subscriptionHandler := handlerSubscription.NewSubscriptionHandler(serviceSubscription.NewSubscriptionService(repositorySubscription.NewSubscriptionRepository(), servicePlan.NewPlanService(repositoryPlan.NewPlanRepository())))
	taskService := background.NewTaskService(serviceSubscription.NewSubscriptionService(repositorySubscription.NewSubscriptionRepository(), servicePlan.NewPlanService(repositoryPlan.NewPlanRepository())))
	return userHandler, planHandler, operatorHandler, subscriptionHandler, taskService
}

func main() {
	router := mux.NewRouter()
	userHandler, planHandler, operatorHandler, subscriptionHandler, taskService := initLayers()
	taskService.StartPeriodicTasks()

	router.HandleFunc("/register-user", userHandler.Create).Methods("POST") // TODO
	router.HandleFunc("/update-user", userHandler.Update).Methods("PUT")    // TODO
	router.HandleFunc("/get-user/{id}", userHandler.Get).Methods("GET")     // TODO

	router.HandleFunc("/get-plan/{id}", planHandler.Get).Methods("GET") // TODO
	router.HandleFunc("/get-plans", planHandler.GetAll).Methods("GET")  // TODO

	router.HandleFunc("/get-operator/{id}", operatorHandler.Get).Methods("GET") // TODO
	router.HandleFunc("/get-operators", operatorHandler.GetAll).Methods("GET")  // TODO

	router.HandleFunc("/add-subscription", subscriptionHandler.AddSubscription).Methods("POST")      // TODO
	router.HandleFunc("/subscription/{id}", subscriptionHandler.Get).Methods("GET")                  // TODO
	router.HandleFunc("/update-key-subscription/{id}", subscriptionHandler.UpdateKey).Methods("PUT") // TODO

	log.Println("Server starting on :8082")
	log.Fatal(http.ListenAndServe(":8082", router))
}
