package validators_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/c0va23/redirector/models"
	"github.com/c0va23/redirector/test/factories"
	"github.com/c0va23/redirector/validators"
)

func TestValidateTarget_EmptyPath(t *testing.T) {
	a := assert.New(t)

	target := factories.
		TargetFactory.
		MustCreateWithOption(map[string]interface{}{
			"Path": "",
		}).(models.Target)

	targetError, valid := validators.ValidateTarget(target)

	a.False(valid)

	a.Equal(
		models.ModelValidationError{
			{
				Name: "path",
				Errors: []models.ValidationError{
					{TranslationKey: "required"},
				},
			},
		},
		targetError,
	)
}

func TestValidateTarget_HTTCodeTooLow(t *testing.T) {
	a := assert.New(t)

	target := factories.
		TargetFactory.
		MustCreateWithOption(map[string]interface{}{
			"HTTPCode": int32(299),
		}).(models.Target)

	targetError, valid := validators.ValidateTarget(target)

	a.False(valid)

	a.Equal(
		models.ModelValidationError{
			{
				Name: "httpCode",
				Errors: []models.ValidationError{
					{TranslationKey: "target.httpCode.outOfRange"},
				},
			},
		},
		targetError,
	)
}

func TestValidateTarget_HTTCodeTooMuch(t *testing.T) {
	a := assert.New(t)

	target := factories.
		TargetFactory.
		MustCreateWithOption(map[string]interface{}{
			"HTTPCode": int32(400),
		}).(models.Target)

	targetError, valid := validators.ValidateTarget(target)

	a.False(valid)

	a.Equal(
		models.ModelValidationError{
			{
				Name: "httpCode",
				Errors: []models.ValidationError{
					{TranslationKey: "target.httpCode.outOfRange"},
				},
			},
		},
		targetError,
	)
}

func TestValidateTarget_HTTCodeMinValue(t *testing.T) {
	a := assert.New(t)

	target := factories.
		TargetFactory.
		MustCreateWithOption(map[string]interface{}{
			"HTTPCode": int32(300),
		}).(models.Target)

	targetError, valid := validators.ValidateTarget(target)

	a.True(valid)

	a.Equal(
		models.ModelValidationError(nil),
		targetError,
	)
}

func TestValidateTarget_HTTCodeMaxValue(t *testing.T) {
	a := assert.New(t)

	target := factories.
		TargetFactory.
		MustCreateWithOption(map[string]interface{}{
			"HTTPCode": int32(399),
		}).(models.Target)

	targetError, valid := validators.ValidateTarget(target)

	a.True(valid)

	a.Equal(
		models.ModelValidationError(nil),
		targetError,
	)
}
