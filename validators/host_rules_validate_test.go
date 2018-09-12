package validators_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/c0va23/redirector/models"
	"github.com/c0va23/redirector/test/factories"
	"github.com/c0va23/redirector/validators"
)

func TestHostRulesValidate_Valid(t *testing.T) {
	a := assert.New(t)

	hostRules := factories.HostRulesFactory.MustCreate().(models.HostRules)

	hostRulesErrors, hostRulesValid := validators.ValidateHostRules(hostRules)

	a.True(hostRulesValid)

	a.Equal(
		models.ModelValidationError(nil),
		hostRulesErrors,
	)
}

func TestHostRulesValidate_EmptyHost(t *testing.T) {
	a := assert.New(t)

	hostRules := factories.
		HostRulesFactory.
		MustCreateWithOption(map[string]interface{}{
			"Host": "",
		}).(models.HostRules)

	hostRulesErrors, hostRulesValid := validators.ValidateHostRules(hostRules)

	a.False(hostRulesValid)

	a.Equal(
		models.ModelValidationError{
			models.FieldValidationError{
				Name: "host",
				Errors: []models.ValidationError{
					{TranslationKey: "errors.required"},
				},
			},
		},
		hostRulesErrors,
	)
}

func TestHostRulesValidate_TooLongHost(t *testing.T) {
	a := assert.New(t)

	longHost := ""

	for i := 0; i <= 256; i++ {
		longHost = longHost + "a"
	}

	hostRules := factories.
		HostRulesFactory.
		MustCreateWithOption(map[string]interface{}{
			"Host": longHost,
		}).(models.HostRules)

	hostRulesErrors, hostRulesValid := validators.ValidateHostRules(hostRules)

	a.False(hostRulesValid)

	a.Equal(
		models.ModelValidationError{
			models.FieldValidationError{
				Name: "host",
				Errors: []models.ValidationError{
					{TranslationKey: "errors.hostRules.host.tooLong"},
				},
			},
		},
		hostRulesErrors,
	)
}

func TestHostRulesValidate_InvalidPatternHost(t *testing.T) {
	a := assert.New(t)

	hostRules := factories.
		HostRulesFactory.
		MustCreateWithOption(map[string]interface{}{
			"Host": "-",
		}).(models.HostRules)

	hostRulesErrors, hostRulesValid := validators.ValidateHostRules(hostRules)

	a.False(hostRulesValid)

	a.Equal(
		models.ModelValidationError{
			models.FieldValidationError{
				Name: "host",
				Errors: []models.ValidationError{
					{TranslationKey: "errors.hostRules.host.invalidPattern"},
				},
			},
		},
		hostRulesErrors,
	)
}

func TestHostRulesValidate_InvalidDefaultTarget(t *testing.T) {
	a := assert.New(t)

	invalidTarget := factories.
		TargetFactory.
		MustCreateWithOption(map[string]interface{}{
			"Path": "",
		}).(models.Target)

	hostRules := factories.
		HostRulesFactory.
		MustCreateWithOption(map[string]interface{}{
			"DefaultTarget": invalidTarget,
		}).(models.HostRules)

	hostRulesErrors, hostRulesValid := validators.ValidateHostRules(hostRules)

	a.False(hostRulesValid)

	a.Equal(
		models.ModelValidationError{
			models.FieldValidationError{
				Name: "defaultTarget.path",
				Errors: []models.ValidationError{
					{TranslationKey: "errors.required"},
				},
			},
		},
		hostRulesErrors,
	)
}

func TestHostRulesValidate_InvalidRule(t *testing.T) {
	a := assert.New(t)

	invalidRule := factories.
		RuleFactory.
		MustCreateWithOption(map[string]interface{}{
			"Resolver": "",
		}).(models.Rule)

	hostRules := factories.
		HostRulesFactory.
		MustCreateWithOption(map[string]interface{}{
			"Rules": []models.Rule{
				invalidRule,
			},
		}).(models.HostRules)

	hostRulesErrors, hostRulesValid := validators.ValidateHostRules(hostRules)

	a.False(hostRulesValid)

	a.Equal(
		models.ModelValidationError{
			models.FieldValidationError{
				Name: "rules.0.resolver",
				Errors: []models.ValidationError{
					{TranslationKey: "errors.rule.resolver.unknown"},
				},
			},
		},
		hostRulesErrors,
	)
}
