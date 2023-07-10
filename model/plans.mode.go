package model

import (
	"time"

	"github.com/google/uuid"
)

type Plan struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	PlanName        string    `gorm:"type:ekyc_schema.plan_type;not null"`
	IsActive        bool      `gorm:"not null"`
	ImageUploadCost float64       `gorm:"not null"`
	FaceMatchCost   float64       `gorm:"not null"`
	OCRCost         float64       `gorm:"not null"`
	DailyBaseCost   float64       `gorm:"not null"`
	CreatedAt       time.Time `gorm:"not null"`
	UpdatedAt       time.Time `gorm:"not null"`
}
