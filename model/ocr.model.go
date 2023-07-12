package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type OCRData struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CustomerID uuid.UUID      `gorm:"type:uuid;not null"`
	ImageID1   uuid.UUID      `gorm:"column:image_id;type:uuid;not null"`
	OCRData    datatypes.JSON `gorm:"column:ocr_data;type:jsonb;not null"`
	CreatedAt  time.Time      `gorm:"not null"`
	UpdatedAt  time.Time      `gorm:"not null"`
}

func (c *OCRData) BeforeCreate(tx *gorm.DB) (err error) {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return
}

func (c *OCRData) BeforeUpdate(tx *gorm.DB) (err error) {
	c.UpdatedAt = time.Now()
	return
}


type OCRAPICalls struct {
    ID             uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
    CustomerID     uuid.UUID  `gorm:"type:uuid;not null"`
    ImageID        uuid.UUID  `gorm:"type:uuid;not null"`
    OCRID          uuid.UUID  `gorm:"type:uuid;not null"`
    APICallCharges float64    `gorm:"not null"`
    CreatedAt      time.Time  `gorm:"not null"`
    UpdatedAt      time.Time  `gorm:"not null"`
}

func (c *OCRAPICalls) BeforeCreate(tx *gorm.DB) (err error) {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return
}

func (c *OCRAPICalls) BeforeUpdate(tx *gorm.DB) (err error) {
	c.UpdatedAt = time.Now()
	return
}