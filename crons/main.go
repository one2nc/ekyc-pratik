package crons

import (
	service "go-ekyc/services"

	"github.com/robfig/cron"
)

func RegisterCron(appService *service.ApplicationService) {
	cron := cron.New()

	scheduleReportCron(cron, appService)

	cron.Start()
}
