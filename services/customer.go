package service

import (
	"errors"
	"go-ekyc/helper"
	"go-ekyc/model"
	"go-ekyc/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ICustomerService interface {
}
type RegisterServiceInput struct {
	PlanName      string
	CustomerEmail string
	CustomerName  string
}

type RegisterCustomerResult struct {
	AccessKey string
	SecretKey string
}

type CustomerService struct {
	customerRepository    repository.ICustomerRepository
	plansRepository       repository.IPlansRepository
	imageRepository       repository.IImageRepository
	faceMatchRepository   repository.IFaceMatchScoreRepository
	ocrRepository         repository.IOCRRepository
	DailyReportRepository repository.IDailyReportsRepository
	RedisRepository       repository.RedisRepository
}

func newCustomerService(customerRepository repository.ICustomerRepository, plansRepository repository.IPlansRepository, imageRepository repository.IImageRepository, faceMatchRepository repository.IFaceMatchScoreRepository, ocrRepository repository.IOCRRepository, dailyReportRepository repository.IDailyReportsRepository, redisRepository repository.RedisRepository) *CustomerService {
	return &CustomerService{
		customerRepository:    customerRepository,
		plansRepository:       plansRepository,
		imageRepository:       imageRepository,
		faceMatchRepository:   faceMatchRepository,
		ocrRepository:         ocrRepository,
		DailyReportRepository: dailyReportRepository,
		RedisRepository:       redisRepository,
	}
}

func (c *CustomerService) RegisterCustomer(serviceInput RegisterServiceInput) (RegisterCustomerResult, error) {
	plan, err := c.plansRepository.FetchPlansByName(serviceInput.PlanName)

	if err != nil {

		return RegisterCustomerResult{}, errors.New("error while fetching plan")
	}

	customer, err := c.customerRepository.GetCustomerByEmail(serviceInput.CustomerEmail)
	if err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {

		return RegisterCustomerResult{}, errors.New("error while fetching customer")

	}
	if customer != (model.Customer{}) {

		return RegisterCustomerResult{}, errors.New("email is already registered")
	}
	accessKey := helper.GenerateRandomString(10)
	secretKey := helper.GenerateRandomString(20)

	customerData := model.Customer{
		Name:      serviceInput.CustomerName,
		Email:     serviceInput.CustomerEmail,
		PlanID:    plan.ID,
		AccessKey: helper.GetMD5Hash(accessKey),
		SecretKey: helper.GetMD5Hash(secretKey),
	}
	err = c.customerRepository.CreateCustomer(&customerData)

	if err != nil {

		return RegisterCustomerResult{}, err
	}

	return RegisterCustomerResult{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}, err
}

func (c *CustomerService) GetCustomerByCredendials(accessKey string, secretKey string) (model.Customer, error) {

	customer, err := c.customerRepository.GetCustomerByCredendials(accessKey, secretKey)
	if err != nil {
		return customer, errors.New("error while fetching customer")

	}
	return customer, nil
}

func (c *CustomerService) CreateCustomerReports(startDate time.Time, endDate time.Time) error {
	customers, err := c.customerRepository.GetCustomersWithPlans()
	if err != nil {
		return err
	}
	// get image upload charges
	imageUploadCharge, err := c.imageRepository.GetImageUploadAPIReport(startDate, endDate)
	if err != nil {
		return err

	}
	// get face-macth charges
	faceMatchApiCharges, err := c.faceMatchRepository.GetFaceMatchAPIReport(startDate, endDate)
	if err != nil {
		return err

	}
	// get ocr charges
	ocrApiCharges, err := c.ocrRepository.GetOCRAPIReport(startDate, endDate)
	if err != nil {
		return err

	}
	// create data for bulk upload
	reports := []model.DailyReport{}
	for _, customer := range customers {
		baseCharge := customer.Plan.DailyBaseCost

		faceMatchCost := faceMatchApiCharges[customer.ID].TotalApiCharge
		faceMatchCount := faceMatchApiCharges[customer.ID].TotalApiCount

		ocrCost := ocrApiCharges[customer.ID].TotalApiCharge
		ocrCount := ocrApiCharges[customer.ID].TotalApiCount

		imageUploadCost := imageUploadCharge[customer.ID].TotalUploadCharges
		imageUploadSize := imageUploadCharge[customer.ID].TotalImageSize
		reports = append(reports, model.DailyReport{
			PlanName:                customer.Plan.PlanName,
			CustomerID:              customer.ID,
			DateOfReport:            startDate,
			DailyBaseCharges:        baseCharge,
			NoOfFaceMatch:           faceMatchCount,
			TotalCostOfFaceMatch:    faceMatchCost,
			TotalImageStorageSizeMb: imageUploadSize,
			TotalImageStorageCost:   imageUploadCost,
			NumberOfOCR:             ocrCount,
			TotalCostOfOCR:          ocrCost,
			TotalAPICallCharges:     ocrCost + faceMatchCost + imageUploadCost,
		})
	}
	err = c.DailyReportRepository.BulkCreateDailyReports(reports)
	if err != nil {
		return err

	}
	return err
}

func (c *CustomerService) GetAggregateReportForCustomer(startDate time.Time, endDate time.Time, customerIds []uuid.UUID) ([]repository.CustomerAggregatedReport, error) {

	reports, err := c.DailyReportRepository.GetCustomersAggregatedReportByDates(startDate, endDate, customerIds)
	for i, report := range reports {
		reports[i].TotalInvoiceAmount = report.TotalBaseCharge + report.TotalAPICallCharges
		reports[i].StartDate = startDate
		reports[i].EndDate = endDate
	}
	return reports, err

}
