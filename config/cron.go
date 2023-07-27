package config

import (
	"log"
	"os"
	"strconv"
)

type CronConfig struct {
	DailyReportExpression string
	ReportLimit           int
}

func NewCronConfig() CronConfig {

	reportLimitString := os.Getenv("REPORT_LIMIT")

	reportLimit, err := strconv.Atoi(reportLimitString)

	if err != nil {
		log.Fatal(err.Error())
	}

	return CronConfig{
		DailyReportExpression: os.Getenv("DAILY_REPORT_CRON_EXPRESSION"),
		ReportLimit:           reportLimit,
	}
}
