package helper

import (
	"go-ekyc/model"
)

func IsImageTypeValid (image_type string) bool{

	return image_type == "face" || image_type == "id_card"
}

func IsImagesComparable (image1 model.Image,image2 model.Image) bool{

	if image1.ImageType == "face" && (image2.ImageType =="face" || image2.ImageType == "id_card"){
		return true
	}
	if image2.ImageType == "face" && (image1.ImageType =="face" || image1.ImageType == "id_card"){
		return true
	}


	return false
	
}