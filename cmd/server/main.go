package main

import (
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"userServer/internal/app"
	yaml "userServer/internal/model/config/YAML"
)

// Middleware для проверки IP
func ipWhitelistMiddleware(allowedIP string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := getClientIP(r)

			if clientIP != allowedIP {
				log.Printf("❌ Доступ запрещен с IP: %s", clientIP)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			log.Printf("✅ Разрешен доступ с IP: %s", clientIP)
			next.ServeHTTP(w, r)
		})
	}
}

// Функция для получения реального IP клиента
func getClientIP(r *http.Request) string {
	// Проверяем заголовки, которые могут быть установлены прокси
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Берем первый IP из списка
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Если заголовков нет, используем RemoteAddr
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

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

	// Применяем middleware ко всем роутам
	router.Use(ipWhitelistMiddleware("85.142.90.53"))

	application.RegisterRoutes(router, cfg)

	addr := ":1288"
	log.Printf("✅✅✅ Сервер запущен на %s ✅✅✅", addr)
	log.Printf("✅ Разрешен только IP: 85.142.90.53")

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("❌❌❌ Ошибка запуска сервера: %v ❌❌❌", err)
	}
}
