package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func FormatErrors(err error) []string {
	var errors []string
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			errorMsg := fieldError.Error()
			if idx := strings.Index(errorMsg, "Error:"); idx != -1 {
				errorMsg = strings.TrimSpace(errorMsg[idx+len("Error:"):])
			}
			errors = append(errors, errorMsg)
		}
	} else {
		errorMsg := err.Error()
		if idx := strings.Index(errorMsg, "Error:"); idx != -1 {
			errorMsg = strings.TrimSpace(errorMsg[idx+len("Error:"):])
		}
		errors = append(errors, errorMsg)
	}
	return errors
}
