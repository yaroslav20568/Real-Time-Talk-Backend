package validator

import (
	"strings"
	"unicode"

	"github.com/gin-gonic/gin/binding"
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

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) == 0 {
		return false
	}

	runes := []rune(password)
	if len(runes) == 0 {
		return false
	}

	if !unicode.IsUpper(runes[0]) {
		return false
	}

	hasSpecialChar := false
	for _, r := range password {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsSpace(r) {
			hasSpecialChar = true
			break
		}
	}

	return hasSpecialChar
}

func Init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("password", validatePassword)
	}
}
