package crons

import (
	"context"
	service "go-ekyc/services"
	"log"
	"time"

	"github.com/robfig/cron"
)

func scheduleReportCron(c *cron.Cron, appService *service.ApplicationService) {

	c.AddFunc("0 0 * * *", func() {

		lockKey := "daily_reports_cron_job_key_1"
		lockValue := true
		expiration := time.Duration.Minutes(1)

		acquired, err := appService.CustomerService.RedisRepository.SetNX(context.Background(), lockKey, lockValue, time.Duration(expiration))

		if err != nil {
			log.Print(err.Error())
		}
		if acquired {

			currentTime := time.Now().UTC()
			startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()-1, 0, 0, 0, 0, time.UTC)
			endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()-1, 23, 59, 59, 0, time.UTC)

			err = appService.CustomerService.CreateCustomerReports(startTime, endTime)

		}

	})

}
