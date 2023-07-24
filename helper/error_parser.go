package helper

import (
	"errors"
	"io"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ErrorParser(err error, object interface{}) []string {

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
		field := e.Field()
		tag := e.Tag()

		structField, ok := reflect.TypeOf(object).Elem().FieldByName(field)
		if !ok {
			continue
		}
		key := ""
		jsonKey, ok := structField.Tag.Lookup("json")
		if ok {
			key = jsonKey
		} else {
			formKey, ok := structField.Tag.Lookup("form")
			if ok {
				key = formKey
			}
		}

		var errorMsg string
		switch tag {
		case "required":
			errorMsg = key + " is required"
		case "oneof":
			errorMsg = key + " should be one of (" + strings.ReplaceAll(e.Param(), " ", "/") + ")"
		case "email":
			errorMsg = key + ": " + "Inavlid email"
		case "uuid":
			errorMsg = key + ": " + "Inavlid uuid"
		default:
			errorMsg = ""
		}

		errorMessagess = append(errorMessagess, errorMsg)
	}

	return errorMessagess
}
