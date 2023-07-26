package service

import "errors"

var (
	ErrPlanNotFound     = errors.New("plan not found")
	ErrCustomerNotFound = errors.New("customer not found")
	ErrImageNotFound    = errors.New("image not found")
	ErrUnknown          = errors.New("unknown error")
	ErrEmailExists      = errors.New("email already registered")
	ErrInvalidImageType = errors.New("invalid image type")
	ErrEmptySlice       = errors.New("empty slice found")
)
