package model

import (
	"time"

	"github.com/google/uuid"
)

type Plan struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	PlanName        string    `gorm:"type:ekyc_schema.plan_type;not null"`
	IsActive        bool      `gorm:"not null"`
	ImageUploadCost int       `gorm:"not null"`
	FaceMatchCost   int       `gorm:"not null"`
	OCRCost         int       `gorm:"not null"`
	DailyBaseCost   int       `gorm:"not null"`
	CreatedAt       time.Time `gorm:"not null"`
	UpdatedAt       time.Time `gorm:"not null"`
}
