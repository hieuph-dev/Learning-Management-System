package validation

import (
	"fmt"
	"lms/src/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitValidation() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("failed to get validation engine")
	}

	RegisterCustomValidation(v)
	return nil
}

func HandlerValidationErrors(err error) gin.H {
	if validationError, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)

		for _, e := range validationError {
			root := strings.Split(e.Namespace(), ".")[0]

			rawPath := strings.TrimPrefix(e.Namespace(), root+".")

			parts := strings.Split(rawPath, ".")

			for i, part := range parts {
				if strings.Contains("part", "[") {
					idx := strings.Index(part, "[")
					base := utils.CamelToSnake(part[:idx])
					index := part[idx:]
					parts[i] = base + index
				} else {
					parts[i] = utils.CamelToSnake(part)
				}
			}

			fieldPath := strings.Join(parts, ".")

			switch e.Tag() {
			case "gt":
				errors[fieldPath] = fmt.Sprintf("%s must be greater than %s", fieldPath, e.Param())
			case "lt":
				errors[fieldPath] = fmt.Sprintf("%s must be less than %s", fieldPath, e.Param())
			case "gte":
				errors[fieldPath] = fmt.Sprintf("%s must be greater than or equal to %s", fieldPath, e.Param())
			case "lte":
				errors[fieldPath] = fmt.Sprintf("%s must be less than or equal to %s", fieldPath, e.Param())
			case "uuid":
				errors[fieldPath] = fmt.Sprintf("%s must be a valid UUID", fieldPath)
			case "slug":
				errors[fieldPath] = fmt.Sprintf("%s can only contain lowercase letters, numbers, hyphens, or periods", fieldPath)
			case "min":
				errors[fieldPath] = fmt.Sprintf("%s must be more than %s characters", fieldPath, e.Param())
			case "max":
				errors[fieldPath] = fmt.Sprintf("%s must be less than %s characters", fieldPath, e.Param())
			case "min_int":
				errors[fieldPath] = fmt.Sprintf("%s must have a value greater than %s", fieldPath, e.Param())
			case "max_int":
				errors[fieldPath] = fmt.Sprintf("%s must have a value less than %s", fieldPath, e.Param())
			case "oneof":
				allowedValues := strings.Join(strings.Split(e.Param(), " "), ",")
				errors[fieldPath] = fmt.Sprintf("%s must be one of the values %s", fieldPath, allowedValues)
			case "required":
				errors[fieldPath] = fmt.Sprintf("%s is required", fieldPath)
			case "search":
				errors[fieldPath] = fmt.Sprintf("%s can only contain lowercase letters, uppercase letters, numbers, and spaces", fieldPath)
			case "email":
				errors[fieldPath] = fmt.Sprintf("%s must be invalid email format", fieldPath)
			case "datetime":
				errors[fieldPath] = fmt.Sprintf("%s must be invalid YYYY-MM-DD format", fieldPath)
			case "email_advanced":
				errors[fieldPath] = fmt.Sprintf("%s is in the prohibited list", fieldPath)
			case "password_strong":
				errors[fieldPath] = fmt.Sprintf("%s must be at least 8 characters, including lowercase letters, uppercase letters, numbers, and special characters", fieldPath)
			case "file_ext":
				allowedValues := strings.Join(strings.Split(e.Param(), " "), ",")
				errors[fieldPath] = fmt.Sprintf("%s only allows files with the following extensions: %s", fieldPath, allowedValues)
			case "course_level":
				errors[fieldPath] = fmt.Sprintf("%s must be one of: beginner, intermediate, advanced", fieldPath)
			case "course_status":
				errors[fieldPath] = fmt.Sprintf("%s must be one of: draft, published, archived", fieldPath)
			case "language_code":
				errors[fieldPath] = fmt.Sprintf("%s must be a valid language code (vi, en)", fieldPath)
			case "positive_float":
				errors[fieldPath] = fmt.Sprintf("%s must be a positive number", fieldPath)
			}
		}

		return gin.H{"error": errors}
	}

	return gin.H{
		"error":  "Invalid request",
		"detail": err.Error(),
	}
}
