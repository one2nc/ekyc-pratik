package helper

import (
	"errors"
	"io"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ErrorParser(err error) []string {

	errorMessagess := []string{}

	if errors.Is(err, io.EOF) {
		errorMessagess = append(errorMessagess, "Empty request body")

		return errorMessagess
	}

	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		errorMessagess = append(errorMessagess, "Something went wrog")

		return errorMessagess
	}
	for _, e := range errs {
		// can translate each error one at a time.
		field := e.Field()
		tag := e.Tag()

		var errorMsg string
		switch tag {
		case "required":
			errorMsg = field + " is required"
		case "oneof":
			errorMsg = field + " should be one of (" + strings.ReplaceAll(e.Param(), " ", "/") + ")"
		case "email":
			errorMsg = "Inavlid email"
		case "uuid":
			errorMsg = "Inavlid uuid"
		default:
			errorMsg = ""
		}

		errorMessagess = append(errorMessagess, errorMsg)
	}

	return errorMessagess
}
