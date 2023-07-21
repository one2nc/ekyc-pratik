package crons

import (
	"go-ekyc/config"
	service "go-ekyc/services"
	"time"

	"github.com/robfig/cron"
)

func scheduleReportCron(c *cron.Cron, config config.CronConfig, appService *service.ApplicationService) {

	c.AddFunc(config.DailyReportExpression, func() {

		currentTime := time.Now().UTC()
		startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()-1, 0, 0, 0, 0, time.UTC)
		endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()-1, 23, 59, 59, 0, time.UTC)
		lockKey := startTime.Unix()

		acquired := appService.CustomerService.DailyReportRepository.AcquireCronLock(lockKey)
		if acquired {

			_ = appService.CustomerService.CreateCustomerReports(startTime, endTime)
			appService.CustomerService.DailyReportRepository.ReleaseCronLock(lockKey)

		}

	})

}
