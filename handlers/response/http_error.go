package response

import (
	"errors"
	service "go-ekyc/services"
	"net/http"
)

func GetHttpStatusAndError(err error) (int, error) {
	switch err {
	case service.ErrCustomerNotFound, service.ErrImageNotFound, service.ErrPlanNotFound:
		return http.StatusNotFound, err
	case service.ErrEmailExists:
		return http.StatusConflict, err
	case service.ErrInvalidImageType:
		return http.StatusUnprocessableEntity, err
	default:

		return http.StatusInternalServerError, errors.New("something went wrong")
	}
}
