package helper

func IsImageTypeValid (image_type string) bool{

	return image_type == "face" || image_type == "id_card"
}