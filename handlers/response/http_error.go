package response

import (
	"errors"
	"fmt"
	service "go-ekyc/services"
	"net/http"
)

func GetHttpStatusAndError(err error) (int, error) {
	fmt.Println("heloo")
	switch err {
	case service.ErrCustomerNotFound, service.ErrImageNotFound, service.ErrPlanNotFound:
		return http.StatusNotFound, err
	case service.ErrEmailExists:
		return http.StatusConflict, err
	case service.ErrInvalidImageType:
		return http.StatusUnprocessableEntity, err
	case service.ErrUnknown:
		fallthrough
	default:

		fmt.Println("hello")
		return http.StatusInternalServerError, errors.New("something went wrong")
	}
}
