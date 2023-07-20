package crons

import (
	"go-ekyc/config"
	service "go-ekyc/services"

	"github.com/robfig/cron"
)

func RegisterCron(appService *service.ApplicationService, config config.CronConfig) {
	cron := cron.New()

	scheduleReportCron(cron, config, appService)

	cron.Start()
}
