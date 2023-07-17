package repository

import (
	"go-ekyc/model"

	"gorm.io/gorm"
)

type IDailyReportsRepository interface {
	BulkCreateDailyReports(reports []model.DailyReport) error
}

type DailyReportsRepository struct {
	dbInstance *gorm.DB
}

func newDailyReportsRepository(db *gorm.DB) IDailyReportsRepository {
	return &DailyReportsRepository{
		dbInstance: db,
	}
}

func (r *DailyReportsRepository) BulkCreateDailyReports(reports []model.DailyReport) error {
	result := r.dbInstance.Create(&reports)

	return result.Error
}
