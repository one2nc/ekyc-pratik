package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DailyReport struct {
	ID                      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CustomerID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();not null"`
	DateOfReport            time.Time `gorm:"not null"`
	DailyBaseCharges        float64   `gorm:"not null"`
	NoOfFaceMatch           int       `gorm:"not null"`
	TotalCostOfFaceMatch    float64   `gorm:"not null"`
	NumberOfOCR             int       `gorm:"not null"`
	TotalCostOfOCR          float64   `gorm:"not null"`
	TotalAPICallCharges     float64   `gorm:"not null"`
	TotalImageStorageSizeMb float64   `gorm:"not null"`
	TotalImageStorageCost   float64   `gorm:"not null"`
	CreatedAt               time.Time `gorm:"not null"`
	UpdatedAt               time.Time `gorm:"not null"`
}

func (d *DailyReport) TableName() string {
	return "ekyc_schema.daily_reports_table"
}
func (d *DailyReport) BeforeCreate(tx *gorm.DB) (err error) {
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()
	return
}

func (d *DailyReport) BeforeUpdate(tx *gorm.DB) (err error) {
	d.UpdatedAt = time.Now()
	return
}
