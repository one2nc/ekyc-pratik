package service_inputs

import (
	"go-ekyc/model"
	"mime/multipart"
)

type UploadImageInput struct {
	Customer  model.Customer
	File      multipart.File
	FileInfo  *multipart.FileHeader
	ImageType string
}
type FaceMatchInput struct {
	Customer model.Customer
	ImageId1 string
	ImageId2 string
}
type OCRInput struct {
	Customer model.Customer
	ImageId string
}
