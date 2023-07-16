package repository

import (
	"go-ekyc/model"

	"github.com/google/uuid"
)

type ImageRepositoryMock struct {
	images []model.Image
}

func (i *ImageRepositoryMock) CreateImage(image *model.Image) error {

	image.ID = uuid.New()
	return nil

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

func newImageMockRepository() IImageRepository {
	return &ImageRepository{}
}
