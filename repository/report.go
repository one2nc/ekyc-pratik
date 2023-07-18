package repository

import (
	"go-ekyc/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerAggregatedReport struct {
	CustomerID              uuid.UUID `json:"customer_id"`
	StartDate               time.Time `json:"start_date_of_report"`
	EndDate                 time.Time `json:"end_date_of_report"`
	TotalBaseCharge         float64   `json:"total_base_charge"`
	TotalFaceMatchCount     int       `json:"total_face_match_count"`
	TotalFaceMatchCost      float64   `json:"total_face_match_cost"`
	TotalOCRCount           int       `json:"total_ocr_count"`
	TotalOCRCost            float64   `json:"total_ocr_cost"`
	TotalImageStorageSizeMb float64   `json:"total_image_storage_size_mb"`
	TotalImageStorageCost   float64   `json:"total_image_storage_cost"`
	TotalAPICallCharges     float64   `json:"total_api_call_charges"`
	TotalInvoiceAmount      float64   `json:"total_invoive_amount"`
	PlanName                string    `json:"plan_name"`
}

type IDailyReportsRepository interface {
	BulkCreateDailyReports(reports []model.DailyReport) error
	GetCustomersAggregatedReportByDates(startDate time.Time, endDate time.Time, customerId []uuid.UUID) ([]CustomerAggregatedReport, error)
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

func (r *DailyReportsRepository) GetCustomersAggregatedReportByDates(startDate time.Time, endDate time.Time, customerIds []uuid.UUID) ([]CustomerAggregatedReport, error) {

	reports := []CustomerAggregatedReport{}

	query := r.dbInstance.Table("ekyc_schema.daily_reports_table").
		Select("daily_reports_table.customer_id,SUM(daily_base_charges) as total_base_charge,SUM(no_of_face_match) as total_face_match_count,SUM(total_cost_of_face_match) as total_face_match_cost,SUM(number_of_ocr) as total_ocr_count, SUM(total_cost_of_ocr) as total_ocr_cost,SUM(total_api_call_charges) as total_api_call_charges, SUM(total_image_storage_size_mb) as total_image_storage_size_mb, SUM(total_image_storage_cost) as total_image_storage_cost,  plans.plan_name").
		Joins(
			"JOIN ekyc_schema.customers ON daily_reports_table.customer_id = customers.id").
		Joins("JOIN ekyc_schema.plans ON customers.plan_id = plans.id")

	if len(customerIds) > 0 {
		query.Where("daily_reports_table.customer_id IN (?) and daily_reports_table.created_at BETWEEN ? AND ?", customerIds, startDate, endDate)
	} else {
		query.Where("daily_reports_table.created_at BETWEEN ? AND ?", startDate, endDate)
	}

	result := query.Group("daily_reports_table.customer_id").Group("plans.plan_name").Scan(&reports)

	return reports, result.Error
}
