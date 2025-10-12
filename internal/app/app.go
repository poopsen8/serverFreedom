package app

import (
	"database/sql"

	"userServer/internal/application/background"
	"userServer/internal/db"
	handlerOperator "userServer/internal/handler/http/operator"
	handlerPlan "userServer/internal/handler/http/plan"
	handlerSubscription "userServer/internal/handler/http/subscription"
	handlerUser "userServer/internal/handler/http/user"
	yaml "userServer/internal/model/config/YAML"
	repoJSONSubscription "userServer/internal/repository/json/subscription"
	repoOperator "userServer/internal/repository/postgres/operetor"
	repoPlan "userServer/internal/repository/postgres/plan"
	repoSubscription "userServer/internal/repository/postgres/subscription"
	repoUser "userServer/internal/repository/postgres/user"
	serviceOperator "userServer/internal/service/operetor"
	servicePlan "userServer/internal/service/plan"
	serviceSubscription "userServer/internal/service/subscription"
	serviceUser "userServer/internal/service/user"

	"github.com/gorilla/mux"
)

type App struct {
	DB                  *sql.DB
	UserHandler         *handlerUser.UserHandler
	PlanHandler         *handlerPlan.PlanHandler
	OperatorHandler     *handlerOperator.OperatorHandler
	SubscriptionHandler *handlerSubscription.SubscriptionHandler
	TaskService         *background.TaskService
}

func New(cfg *yaml.RouteConfig) (*App, error) {
	database, err := db.NewPostgres(&cfg.DB)
	if err != nil {
		return nil, err
	}

	operatorRepo := repoOperator.NewOperetorRepository(database)
	userRepo := repoUser.NewUserRepository(database)
	planRepo := repoPlan.NewPlanRepository(database)
	subscriptionRepo := repoSubscription.NewSubscriptionRepository(database)
	subscriptionRepo2 := repoJSONSubscription.NewSubscriptionRepository()

	operatorService := serviceOperator.NewOperatorService(operatorRepo)
	planService := servicePlan.NewPlanService(planRepo)
	userService := serviceUser.NewUserService(userRepo, *operatorService)
	subscriptionService := serviceSubscription.NewSubscriptionService(
		subscriptionRepo,
		subscriptionRepo2,
		*planService,
		*userService,
		*operatorService,
	)

	operatorHandler := handlerOperator.NewOperatorHandler(operatorService, *cfg)
	userHandler := handlerUser.NewUserHandler(userService, *cfg)
	planHandler := handlerPlan.NewPlanHandler(planService, *cfg)
	subscriptionHandler := handlerSubscription.NewSubscriptionHandler(subscriptionService, *cfg)

	taskService := background.NewTaskService(subscriptionService)

	return &App{
		DB:                  database,
		UserHandler:         userHandler,
		PlanHandler:         planHandler,
		OperatorHandler:     operatorHandler,
		SubscriptionHandler: subscriptionHandler,
		TaskService:         taskService,
	}, nil
}

func (a *App) RegisterRoutes(router *mux.Router) {
	// --- USER ROUTES ---
	router.HandleFunc("/register-user", a.UserHandler.Create).Methods("POST") // принимает id, username
	router.HandleFunc("/update-user", a.UserHandler.Update).Methods("PUT")    // принимает одно значение на изменение и id пользователя - MobileOperatorID, IsTrial
	router.HandleFunc("/user/{id}", a.UserHandler.User).Methods("GET")        // принимает id

	router.HandleFunc("/users", a.UserHandler.Users).Methods("GET")

	// --- PLAN ROUTES ---
	router.HandleFunc("/plan/{id}", a.PlanHandler.Plan).Methods("GET") // принимает id
	router.HandleFunc("/plans", a.PlanHandler.Plans).Methods("GET")    // получить все планы

	// --- OPERATOR ROUTES ---
	router.HandleFunc("/operator/{id}", a.OperatorHandler.Operator).Methods("GET") // принимает id
	router.HandleFunc("/operators", a.OperatorHandler.Operators).Methods("GET")    // получить всех операторов

	// --- SUBSCRIPTION ROUTES ---
	router.HandleFunc("/add-subscription", a.SubscriptionHandler.AddSubscription).Methods("POST")      // принимает user_id, plan_id, create_at, expires_at
	router.HandleFunc("/subscription/{id}", a.SubscriptionHandler.Subscription).Methods("GET")         // принимает user_id
	router.HandleFunc("/update-key-subscription/{id}", a.SubscriptionHandler.UpdateKey).Methods("PUT") // принимает user_id

	router.HandleFunc("/subscriptions", a.SubscriptionHandler.Subscriptions).Methods("GET") // принимает user_id

}

func (a *App) Close() {
	if a.DB != nil {
		_ = a.DB.Close()
	}
}
