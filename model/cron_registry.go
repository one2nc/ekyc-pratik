package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Your GORM model representing the table
type CronRegistry struct {
	ID               uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name             string         `gorm:"type:varchar;not null"`
	UniqueIdentifier string         `gorm:"type:varchar;unique;not null"`
	Metadata         datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt        time.Time      `gorm:"not null"`
	UpdatedAt        time.Time      `gorm:"not null"`
}

func (CronRegistry) TableName() string {
	return "ekyc_schema.cron_registry"
}
func (c *CronRegistry) BeforeCreate(tx *gorm.DB) (err error) {
	c.CreatedAt = time.Now().UTC()
	c.UpdatedAt = time.Now().UTC()
	return
}

func (c *CronRegistry) BeforeUpdate(tx *gorm.DB) (err error) {
	c.UpdatedAt = time.Now().UTC()
	return
}
