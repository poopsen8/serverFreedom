package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"userServer/internal/app"
	yaml "userServer/internal/model/config/YAML"
)

func main() {
	cfg, err := yaml.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("❌❌❌ Ошибка загрузки конфигурации: %v ❌❌❌", err)
	}

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("❌❌❌ Ошибка инициализации приложения: %v ❌❌❌", err)
	}
	application.TaskService.StartPeriodicTasks()

	router := mux.NewRouter()
	application.RegisterRoutes(router, cfg)

	addr := ":1288"
	log.Printf("✅✅✅ Сервер запущен на %s ✅✅✅", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("❌❌❌ Ошибка запуска сервера: %v ❌❌❌", err)
	}
}
