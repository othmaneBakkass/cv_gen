package tone

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	apperror "github.com/othmaneBakkass/cv_gen/internal/common/appError"
)

type JSONSchema struct {
	Template  string             `json:"template" validate:"required"`
	FileName  string             `json:"fileName" validate:"required"`
	Head      HeadSchema         `json:"head" validate:"required"`
	Education []EducationSchema  `json:"education" validate:"required"`
	Jobs      []JobSchema        `json:"jobs" validate:"required"`
	Languages []LanguageSchema   `json:"languages" validate:"required"`
}

type HeadSchema struct {
	FullName string `json:"fullName" validate:"required"`
	Address  string `json:"address" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Email    string `json:"email" validate:"required"`
}

type EducationSchema struct {
	School      string `json:"school" validate:"required"`
	Location    string `json:"location" validate:"required"`
	StartedAt   string `json:"startedAt" validate:"required"`
	EndedAt     string `json:"endedAt" validate:"required"`
	Degree      string `json:"degree" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type JobSchema struct {
	Company    string   `json:"company" validate:"required"`
	Location   string   `json:"location" validate:"required"`
	Position   string   `json:"position" validate:"required"`
	StartedAt  string   `json:"startedAt" validate:"required"`
	EndedAt    string   `json:"endedAt" validate:"required"`
	Tools      []string `json:"tools" validate:"required,dive,required"`
	Highlights []string `json:"highlights" validate:"required,dive,required"`
}

type LanguageSchema struct {
	Language string `json:"language" validate:"required"`
	Level    string `json:"level" validate:"required"`
}

func Validate(schema *JSONSchema) (*JSONSchema, error) {
	v := validator.New()

	if err := v.Struct(schema); err != nil {
		// Handle internal validator misconfiguration
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return nil, apperror.New(
				"Internal validation failure",
				"An unexpected internal error occurred while validating the JSON file. Please ensure the JSON follows the latest schema.",
				apperror.ErrorCodeArgs,
				apperror.ErrorSensitivityPublic,
				[]apperror.AppErrorIssue{},
			)
		}

		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var issues []apperror.AppErrorIssue
			for _, ve := range validationErrors {
				// Create user-friendly issue without exposing values or internal tags
				issues = append(issues, apperror.AppErrorIssue{
					Title:       fmt.Sprintf("Invalid field: %s", ve.Field()),
					Detail:      fmt.Sprintf("The field '%s' is missing or invalid. Please check the latest JSON schema.", ve.Field()),
					Sensitivity: apperror.ErrorSensitivityPublic,
				})
			}

			return nil, apperror.New(
				"Invalid JSON format",
				"The JSON file does not match the expected format. Please check the latest JSON schema.",
				apperror.ErrorCodeArgs,
				apperror.ErrorSensitivityPublic,
				issues,
			)
		}
	}

	return schema, nil
}
