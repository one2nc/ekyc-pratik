package config

import "os"

type CronConfig struct {
	DailyReportExpression string
}

func NewCronConfig() CronConfig {
	return CronConfig{
		DailyReportExpression: os.Getenv("DAILY_REPORT_CRON_EXPRESSION"),
	}
}
