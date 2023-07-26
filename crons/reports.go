package crons

import (
	"go-ekyc/config"
	service "go-ekyc/services"
	"time"

	"github.com/robfig/cron"
)

func scheduleReportCron(c *cron.Cron, config config.CronConfig, appService *service.ApplicationService) (error){

err :=	c.AddFunc(config.DailyReportExpression, func() {
		currentTime := time.Now().UTC()
		startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()-1, 0, 0, 0, 0, time.UTC)
		endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()-1, 23, 59, 59, 0, time.UTC)

			_ = appService.CustomerService.CreateCustomerReportsCron(startTime, endTime)


	})

	if err != nil {
		return err
	}
return nil
}
