package validators

import (
	"github.com/c0va23/redirector/models"
)

func addFieldError(
	modelError models.ModelValidationError,
	fieldName string,
	translationKey string,
) models.ModelValidationError {
	for _, fieldError := range modelError {
		if fieldName == fieldError.Name {
			fieldError.Errors = append(
				fieldError.Errors,
				models.ValidationError{
					TranslationKey: translationKey,
				},
			)
			return modelError
		}
	}

	modelError = append(
		modelError,
		models.FieldValidationError{
			Name: fieldName,
			Errors: []models.ValidationError{
				{TranslationKey: translationKey},
			},
		},
	)
	return modelError
}

func isEmptyModelError(modelError models.ModelValidationError) bool {
	return 0 == len(modelError)
}

func addEmbedError(
	modelError models.ModelValidationError,
	embedName string,
	embedModelError models.ModelValidationError,
) models.ModelValidationError {
	for _, embedFieldError := range embedModelError {
		modelError = append(
			modelError,
			models.FieldValidationError{
				Name:   embedName + "." + embedFieldError.Name,
				Errors: embedFieldError.Errors,
			},
		)
	}
	return modelError
}
