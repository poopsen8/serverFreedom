// internal/application/background/task_service.go
package background

import (
	"time"
)

type service interface {
	BackgroundCheck() error
}

type TaskService struct {
	subscriptionService service
}

func NewTaskService(subscriptionService service) *TaskService {
	return &TaskService{
		subscriptionService: subscriptionService,
	}
}

func (ts *TaskService) StartPeriodicTasks() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			ts.subscriptionService.BackgroundCheck()
		}
	}()
}
