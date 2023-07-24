package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Image struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CustomerID    uuid.UUID `gorm:"type:uuid;not null"`
	FilePath      string    `gorm:"not null"`
	FileExtension string    `gorm:"not null"`
	FileSizeMB    float64   `gorm:"not null"`
	ImageType     string    `gorm:"not null"`
	CreatedAt     time.Time `gorm:"not null"`
	UpdatedAt     time.Time `gorm:"not null"`
}

func (c *Image) BeforeCreate(tx *gorm.DB) (err error) {
	c.CreatedAt = time.Now().UTC()
	c.UpdatedAt = time.Now().UTC()
	return
}

func (c *Image) BeforeUpdate(tx *gorm.DB) (err error) {
	c.UpdatedAt = time.Now().UTC()
	return
}

type ImageUploadAPICall struct {
	ID                  uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CustomerID          uuid.UUID `gorm:"type:uuid;not null"`
	ImageID             uuid.UUID `gorm:"type:uuid;not null"`
	ImageStorageCharges float64   `gorm:"not null"`
	CreatedAt           time.Time `gorm:"not null"`
	UpdatedAt           time.Time `gorm:"not null"`
}

func (i *ImageUploadAPICall) BeforeCreate(tx *gorm.DB) error {
	i.CreatedAt = time.Now().UTC()
	i.UpdatedAt = time.Now().UTC()
	return nil
}

func (i *ImageUploadAPICall) BeforeUpdate(tx *gorm.DB) error {
	i.UpdatedAt = time.Now().UTC()
	return nil
}
