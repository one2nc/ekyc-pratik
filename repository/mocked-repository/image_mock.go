package mockedrepository

import (
	"go-ekyc/model"
	"go-ekyc/repository"
	"time"

	"github.com/google/uuid"
)

type ImageRepositoryMock struct {
	images                []model.Image
	imageUploadApiReports map[uuid.UUID]repository.ImageUploadAPIReport
}

func newImageMockRepository() repository.IImageRepository {
	return &ImageRepositoryMock{}
}

func (i *ImageRepositoryMock) CreateImage(image *model.Image) error {

	image.ID = uuid.New()
	return nil

}
func (i *ImageRepositoryMock) GetImageUploadAPIReport(startDate time.Time, endDate time.Time) (map[uuid.UUID]repository.ImageUploadAPIReport, error) {

	return i.imageUploadApiReports, nil
}

func (i *ImageRepositoryMock) CreateImageUploadRecord(imageUploadData *model.ImageUploadAPICall) error {

	imageUploadData.ID = uuid.New()

	return nil
}
func (i *ImageRepositoryMock) FindImagesByIdForCustomer(imageIds []string, customerId string) ([]model.Image, error) {

	images := []model.Image{}

	for _, image := range images {

		for _, imageId := range imageIds {
			if imageId == image.ID.String() {
				images = append(images, image)
			}
		}
	}
	return images, nil
}