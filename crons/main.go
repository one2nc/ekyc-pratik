package crons

import (
	"go-ekyc/config"
	service "go-ekyc/services"
	"log"

	"github.com/robfig/cron"
)

func RegisterCron(appService *service.ApplicationService, config config.CronConfig) {
	cron := cron.New()

	err := scheduleReportCron(cron, config, appService)

	if err != nil {
		log.Fatal(err.Error())
	}

	cron.Start()
}
