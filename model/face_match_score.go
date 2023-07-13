package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FaceMatchScore struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CustomerID uuid.UUID `gorm:"type:uuid;not null"`
	ImageID1   uuid.UUID `gorm:"column:image_id_1;type:uuid;not null"`
	ImageID2   uuid.UUID `gorm:"column:image_id_2;type:uuid;not null"`
	Score      int       `gorm:"not null"`
	CreatedAt  time.Time `gorm:"not null"`
	UpdatedAt  time.Time `gorm:"not null"`
}

// Define the foreign key relationships
func (FaceMatchScore) TableName() string {
	return "ekyc_schema.face_match_score"
}

func (fms *FaceMatchScore) BeforeCreate(tx *gorm.DB) error {
	fms.CreatedAt = time.Now()
	fms.UpdatedAt = time.Now()
	return nil
}

func (fms *FaceMatchScore) BeforeUpdate(tx *gorm.DB) error {
	fms.UpdatedAt = time.Now()
	return nil
}

type FaceMatchAPICall struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();not null;primaryKey"`
	CustomerID     uuid.UUID `gorm:"type:uuid;not null"`
	ScoreID        uuid.UUID `gorm:"type:uuid;not null"`
	APICallCharges float64   `gorm:"column:api_call_charges;not null"`
	CreatedAt      time.Time `gorm:"not null"`
	UpdatedAt      time.Time `gorm:"not null"`
}

// Define the foreign key relationships
func (FaceMatchAPICall) TableName() string {
	return "ekyc_schema.face_match_api_calls"
}

func (fm *FaceMatchAPICall) BeforeCreate(tx *gorm.DB) error {
	fm.CreatedAt = time.Now()
	fm.UpdatedAt = time.Now()
	return nil
}

func (fm *FaceMatchAPICall) BeforeUpdate(tx *gorm.DB) error {
	fm.UpdatedAt = time.Now()
	return nil
}
