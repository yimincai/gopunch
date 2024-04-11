package cronjob

import (
	"github.com/robfig/cron/v3"
	"github.com/yimincai/gopunch/internal/service"
	"github.com/yimincai/gopunch/pkg/logger"
)

type Cron struct {
	c   *cron.Cron
	svc service.Service
}

// Initialize cronjob add schedule jobs to cronjob
func (cj *Cron) initialize() {
	_, err := cj.c.AddFunc("30 7 * * *", func() {
		err := cj.svc.DefaultSchedulePunchAllUsers()
		if err != nil {
			logger.Error(err)
		}
		logger.Info("All Users Punch Done")
	})
	if err != nil {
		panic(err)
	}

	_, err = cj.c.AddFunc("0 18 * * *", func() {
		err := cj.svc.DefaultSchedulePunchAllUsers()
		if err != nil {
			logger.Error(err)
		}
		logger.Info("All Users Punch Done")
	})
	if err != nil {
		panic(err)
	}
}

// Start cronjob in goroutine
func (cj *Cron) Start() {
	cj.initialize()
	cj.c.Start()
}

// Stop cronjob
func (cj *Cron) Stop() {
	cj.c.Stop()
}

func New(
	logger *logger.Logger,
	svc service.Service,
) *Cron {
	return &Cron{
		c:   cron.New(cron.WithLogger(logger), cron.WithChain(cron.Recover(logger))),
		svc: svc,
	}
}
