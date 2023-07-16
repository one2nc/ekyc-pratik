package requests

type FaceMatchRequest struct {
	ImageId1 string `json:"image_id_1" binding:"required,uuid"`
	ImageId2 string `json:"image_id_2" binding:"required,uuid"`
}

type OCRRequest struct {
	ImageId1 string `json:"image_id" binding:"required,uuid"`
}
