package repository

import (
	"go-ekyc/model"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ICronRegistryRepository interface {
	CreateCronRecordNX(cronData *model.CronRegistry) (bool, error)
	GetCronByUniqueIdTX(tx *gorm.DB, uniqueID string) (model.CronRegistry, error)
	BeginTX() *gorm.DB
	CommitTX(tx *gorm.DB) *gorm.DB
	RollbackTx(tx *gorm.DB) *gorm.DB
	UpdateCronMetadataByUniqueIdTX(tx *gorm.DB, uniqueID string, metadata datatypes.JSON) (model.CronRegistry, error)
}
type CronRegistryRepository struct {
	dbInstance *gorm.DB
}

func newCronRegistryRepository(db *gorm.DB) ICronRegistryRepository {
	return &CronRegistryRepository{
		dbInstance: db,
	}
}

func (c *CronRegistryRepository) BeginTX() *gorm.DB {

	return c.dbInstance.Begin()
}
func (c *CronRegistryRepository) CommitTX(tx *gorm.DB) *gorm.DB {

	return tx.Commit()
}
func (c *CronRegistryRepository) RollbackTx(tx *gorm.DB) *gorm.DB {

	return tx.Rollback()
}

func (c *CronRegistryRepository) GetCronByUniqueIdTX(tx *gorm.DB, uniqueID string) (model.CronRegistry, error) {
	cron := model.CronRegistry{}
	result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&cron, "unique_identifier = ?", uniqueID)

	return cron, result.Error
}
func (c *CronRegistryRepository) UpdateCronMetadataByUniqueIdTX(tx *gorm.DB, uniqueID string, metadata datatypes.JSON) (model.CronRegistry, error) {
	cron := model.CronRegistry{}
	result := tx.Model(&cron).Where("unique_identifier = ?", uniqueID).Update("metadata", metadata)

	return cron, result.Error
}
func (c *CronRegistryRepository) CreateCronRecordNX(cronData *model.CronRegistry) (bool, error) {

	query := "INSERT INTO ekyc_schema.cron_registry (name,unique_identifier, metadata, created_at,updated_at) VALUES (?, ?, ?, ?, ?) ON CONFLICT DO NOTHING"
	result := c.dbInstance.Exec(query, cronData.Name, cronData.UniqueIdentifier, cronData.Metadata, time.Now().UTC(), time.Now().UTC())

	if result.Error != nil {
		return false, result.Error
	}

	rowsAffected := result.RowsAffected

	return rowsAffected > 0, nil
}
