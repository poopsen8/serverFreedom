package main

import (
	"log"
	"net/http"
	"userServer/internal/application/background"
	handlerOperator "userServer/internal/handler/http/operator"
	handlerSubscription "userServer/internal/handler/http/subscription"

	handlersPlan "userServer/internal/handler/http/plan"
	userHandler "userServer/internal/handler/http/user"
	yaml "userServer/internal/model/config/YAML"
	"userServer/internal/repository/json/subscription"
	repositoryOperetor "userServer/internal/repository/postgres/operetor"
	repositoryPlan "userServer/internal/repository/postgres/plan"
	repositorySubscription "userServer/internal/repository/postgres/subscription"
	userRepository "userServer/internal/repository/postgres/user"
	serviceOperator "userServer/internal/service/operetor"
	servicePlan "userServer/internal/service/plan"
	serviceSubscription "userServer/internal/service/subscription"
	userService "userServer/internal/service/user"

	"github.com/gorilla/mux"
)

func initLayers(routesCfg yaml.RouteConfig) (
	*userHandler.UserHandler,
	*handlersPlan.PlanHandler,
	*handlerOperator.OperatorHandler,
	*handlerSubscription.SubscriptionHandler,
	*background.TaskService,
) {
	operatorRepo := repositoryOperetor.NewOperetorRepository()
	userRepo := userRepository.NewUserRepository()
	planRepo := repositoryPlan.NewPlanRepository()
	subscriptionRepo := repositorySubscription.NewSubscriptionRepository()
	subscriptionRepo2 := subscription.NewSubscriptionRepository()

	operatorService := serviceOperator.NewOperatorService(operatorRepo)
	planService := servicePlan.NewPlanService(planRepo)
	userService := userService.NewUserService(userRepo, *operatorService)
	subscriptionService := serviceSubscription.NewSubscriptionService(
		subscriptionRepo,
		subscriptionRepo2,
		*planService,
		*userService,
		*operatorService,
	)

	operatorHandler := handlerOperator.NewOperatorHandler(operatorService, routesCfg)
	userHandler := userHandler.NewUserHandler(userService, routesCfg)
	planHandler := handlersPlan.NewPlanHandler(planService, routesCfg)
	subscriptionHandler := handlerSubscription.NewSubscriptionHandler(subscriptionService, routesCfg)
	taskService := background.NewTaskService(subscriptionService)

	return userHandler, planHandler, operatorHandler, subscriptionHandler, taskService
}

func main() {

	routesCfg, err := yaml.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	userHandler, planHandler, operatorHandler, subscriptionHandler, taskService := initLayers(*routesCfg)
	taskService.StartPeriodicTasks()

	userHandler.RegisterRoutes(router)
	planHandler.RegisterRoutes(router)
	operatorHandler.RegisterRoutes(router)
	subscriptionHandler.RegisterRoutes(router)

	log.Println("Server starting on :8082")
	log.Fatal(http.ListenAndServe(":8082", router))
}
