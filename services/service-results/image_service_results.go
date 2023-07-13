package service_results

import "github.com/google/uuid"


type ImageUploadResult struct {
	ImageId uuid.UUID
}
type FaceMatchResult struct {
	Score int
}
type OCRResult struct {
	Data interface{}
}