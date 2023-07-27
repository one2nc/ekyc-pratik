package service

import (
	"encoding/json"
	"go-ekyc/helper"
	"go-ekyc/model"
	"go-ekyc/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
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
	customerRepository     repository.ICustomerRepository
	plansRepository        repository.IPlansRepository
	imageRepository        repository.IImageRepository
	faceMatchRepository    repository.IFaceMatchScoreRepository
	ocrRepository          repository.IOCRRepository
	DailyReportRepository  repository.IDailyReportsRepository
	RedisRepository        repository.RedisRepository
	CronRegistryRepository repository.ICronRegistryRepository
}

func newCustomerService(customerRepository repository.ICustomerRepository, plansRepository repository.IPlansRepository, imageRepository repository.IImageRepository, faceMatchRepository repository.IFaceMatchScoreRepository, ocrRepository repository.IOCRRepository, dailyReportRepository repository.IDailyReportsRepository, redisRepository repository.RedisRepository, cronRegistryRepository repository.ICronRegistryRepository) *CustomerService {
	return &CustomerService{
		customerRepository:     customerRepository,
		plansRepository:        plansRepository,
		imageRepository:        imageRepository,
		faceMatchRepository:    faceMatchRepository,
		ocrRepository:          ocrRepository,
		DailyReportRepository:  dailyReportRepository,
		RedisRepository:        redisRepository,
		CronRegistryRepository: cronRegistryRepository,
	}
}

func (c *CustomerService) RegisterCustomer(serviceInput RegisterServiceInput) (RegisterCustomerResult, error) {
	plan, err := c.plansRepository.FetchPlansByName(serviceInput.PlanName)

	if err != nil {

		return RegisterCustomerResult{}, ErrPlanNotFound
	}

	customer, err := c.customerRepository.GetCustomerByEmail(serviceInput.CustomerEmail)
	if err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {

		return RegisterCustomerResult{}, ErrUnknown

	}
	if customer != (model.Customer{}) {

		return RegisterCustomerResult{}, ErrEmailExists
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

		return RegisterCustomerResult{}, ErrUnknown
	}

	return RegisterCustomerResult{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}, nil
}

func (c *CustomerService) GetCustomerByCredendials(accessKey string, secretKey string) (model.Customer, error) {

	customer, err := c.customerRepository.GetCustomerByCredendials(accessKey, secretKey)
	if err != nil {

		if err == gorm.ErrRecordNotFound {
			return customer, ErrCustomerNotFound
		}
		return customer, ErrUnknown
	}
	return customer, nil
}

func (c *CustomerService) CreateCustomerReportsCron(startDate time.Time, endDate time.Time, limit int) error {

	name := "DAILY_REPORTS_CRON"
	uniqueIdentifer := name + "_" + startDate.String()
	createCronMetaData := repository.DailyReportCronMetaData{
		LastOffset: 0,
	}
	jsonData, err := json.Marshal(&createCronMetaData)
	if err != nil {
		return err
	}
	data := model.CronRegistry{
		Name:             name,
		UniqueIdentifier: uniqueIdentifer,
		Metadata:         datatypes.JSON(jsonData),
	}
	_, err = c.CronRegistryRepository.CreateCronRecordNX(&data)

	numOfRetriesToFetchCronRecord := 1
	maxNumOfRetriesToFetchCronRecord := 5
	for {

		tx := c.CronRegistryRepository.BeginTX()
		if tx.Error != nil {
			continue
		}

		cron, err := c.CronRegistryRepository.GetCronByUniqueIdTX(tx, uniqueIdentifer)
		if err != nil {
			c.CronRegistryRepository.RollbackTx(tx)
			if err != gorm.ErrRecordNotFound {
				return err
			} else {
				if numOfRetriesToFetchCronRecord > maxNumOfRetriesToFetchCronRecord {
					return err
				}
				numOfRetriesToFetchCronRecord++
			}
			continue
		}

		metadata := &repository.DailyReportCronMetaData{}
		err = json.Unmarshal(cron.Metadata, &metadata)
		if err != nil {
			c.CronRegistryRepository.RollbackTx(tx)
			return err
		}

		offset := metadata.LastOffset
		metadata.LastOffset += limit
		updateMetadataJson, err := json.Marshal(metadata)
		if err != nil {
			c.CronRegistryRepository.RollbackTx(tx)
			continue
		}

		_, err = c.CronRegistryRepository.UpdateCronMetadataByUniqueIdTX(tx, cron.UniqueIdentifier, datatypes.JSON(string(updateMetadataJson)))
		if err != nil {
			c.CronRegistryRepository.RollbackTx(tx)
			continue
		}

		result := c.CronRegistryRepository.CommitTX(tx)

		if result.Error != nil {
			continue
		}

		numberOFRetries := 5

		for i := 0; i < numberOFRetries; i++ {
			err := c.CreateCustomerReports(startDate, endDate, limit, offset)

			if err == nil {
				break
			}

			if err == ErrEmptySlice {
				return nil
			}

		}

	}

	// return nil
}
func (c *CustomerService) CreateCustomerReports(startDate time.Time, endDate time.Time, limit int, offset int) error {
	customers, err := c.customerRepository.GetCustomersWithPlans(limit, offset, endDate)
	if err != nil {
		return ErrUnknown
	}

	customerIds := []uuid.UUID{}
	for _, customer := range customers {
		customerIds = append(customerIds, customer.ID)
	}
	if len(customerIds) == 0 {
		return ErrEmptySlice
	}
	// get image upload charges
	imageUploadCharge, err := c.imageRepository.GetImageUploadAPIReport(startDate, endDate, customerIds)
	if err != nil {

		return ErrUnknown

	}
	// get face-macth charges
	faceMatchApiCharges, err := c.faceMatchRepository.GetFaceMatchAPIReport(startDate, endDate, customerIds)
	if err != nil {

		return ErrUnknown

	}
	// get ocr charges
	ocrApiCharges, err := c.ocrRepository.GetOCRAPIReport(startDate, endDate, customerIds)
	if err != nil {
		return ErrUnknown

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

		return ErrUnknown

	}

	return nil
}

func (c *CustomerService) GetAggregateReportForCustomer(startDate time.Time, endDate time.Time, customerIds []uuid.UUID) ([]repository.CustomerAggregatedReport, error) {

	reports, err := c.DailyReportRepository.GetCustomersAggregatedReportByDates(startDate, endDate, customerIds)
	for i, report := range reports {
		reports[i].TotalInvoiceAmount = report.TotalBaseCharge + report.TotalAPICallCharges
		reports[i].StartDate = startDate
		reports[i].EndDate = endDate
	}

	if err != nil {
		return reports, ErrUnknown
	}
	return reports, nil

}
