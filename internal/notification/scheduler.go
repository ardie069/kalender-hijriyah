package notification

import (
	"context"
	"github.com/robfig/cron/v3"
)

// Scheduler handles cron scheduling for notifications
type Scheduler struct {
	cron    *cron.Cron
	service *Service
}

func NewScheduler(service *Service) *Scheduler {
	return &Scheduler{
		cron:    cron.New(),
		service: service,
	}
}

func (s *Scheduler) Start() {
	s.cron.AddFunc("0 8 * * *", func() {
		s.service.GenerateNotifications(context.Background())
	})
	s.cron.Start()
}
