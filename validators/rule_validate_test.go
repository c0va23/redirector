package validators_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/c0va23/redirector/models"
	"github.com/c0va23/redirector/test/factories"
	"github.com/c0va23/redirector/validators"
)

func TestValidateRule_EmptyRule(t *testing.T) {
	a := assert.New(t)

	ruleError, valid := validators.ValidateRule(models.Rule{})
	a.False(valid)
	a.Equal(
		models.ModelValidationError{
			{
				Name: "resolver",
				Errors: []models.ValidationError{
					{TranslationKey: "rule.resolver.unknown"},
				},
			},
			{
				Name: "sourcePath",
				Errors: []models.ValidationError{
					{TranslationKey: "required"},
				},
			},
			{
				Name: "target.httpCode",
				Errors: []models.ValidationError{
					{TranslationKey: "target.httpCode.outOfRange"},
				},
			},
			{
				Name: "target.path",
				Errors: []models.ValidationError{
					{TranslationKey: "required"},
				},
			},
		},
		ruleError,
	)
}

func TestValidateRule_SimpleResolver(t *testing.T) {
	a := assert.New(t)

	rule, _ := factories.RuleFactory.MustCreateWithOption(map[string]interface{}{
		"Resolver": models.RuleResolverSimple,
	}).(models.Rule)

	ruleError, valid := validators.ValidateRule(rule)

	a.True(valid)
	a.Equal(
		models.ModelValidationError(nil),
		ruleError,
	)
}

func TestValidateRule_ValidPatternWithValidTarget(t *testing.T) {
	a := assert.New(t)

	target := factories.
		TargetFactory.
		MustCreateWithOption(map[string]interface{}{
			"Path": "https://example.com/posts/{0}",
		}).(models.Target)

	rule, _ := factories.RuleFactory.MustCreateWithOption(map[string]interface{}{
		"Resolver":   models.RuleResolverPattern,
		"SourcePath": "/r/(\\d+)",
		"Target":     target,
	}).(models.Rule)

	ruleError, valid := validators.ValidateRule(rule)

	a.True(valid)
	a.Equal(
		models.ModelValidationError(nil),
		ruleError,
	)
}

func TestValidateRule_ValidPatternWithUnorderedPlaceholders(t *testing.T) {
	a := assert.New(t)

	target := factories.
		TargetFactory.
		MustCreateWithOption(map[string]interface{}{
			"Path": "https://example.com/groups/{1}/posts/{0}",
		}).(models.Target)

	rule, _ := factories.RuleFactory.MustCreateWithOption(map[string]interface{}{
		"Resolver":   models.RuleResolverPattern,
		"SourcePath": "/r/(\\d+)-(\\d+)",
		"Target":     target,
	}).(models.Rule)

	ruleError, valid := validators.ValidateRule(rule)

	a.True(valid)
	a.Equal(
		models.ModelValidationError(nil),
		ruleError,
	)
}

func TestValidateRule_ValidPatternWithoutPlaceholder(t *testing.T) {
	a := assert.New(t)

	target := factories.
		TargetFactory.
		MustCreateWithOption(map[string]interface{}{
			"Path": "https://example.com/posts/",
		}).(models.Target)

	rule, _ := factories.RuleFactory.MustCreateWithOption(map[string]interface{}{
		"Resolver":   models.RuleResolverPattern,
		"SourcePath": "/r/(\\d+)",
		"Target":     target,
	}).(models.Rule)

	ruleError, valid := validators.ValidateRule(rule)

	a.False(valid)
	a.Equal(
		models.ModelValidationError{
			{
				Name: "target.path",
				Errors: []models.ValidationError{
					{TranslationKey: "target.path.missedPlaceholder"},
				},
			},
		},
		ruleError,
	)
}

func TestValidateRule_ValidPatternWithInvalidPlaceholder(t *testing.T) {
	a := assert.New(t)

	target := factories.
		TargetFactory.
		MustCreateWithOption(map[string]interface{}{
			"Path": "https://example.com/posts/{1}",
		}).(models.Target)

	rule, _ := factories.RuleFactory.MustCreateWithOption(map[string]interface{}{
		"Resolver":   models.RuleResolverPattern,
		"SourcePath": "/r/(\\d+)",
		"Target":     target,
	}).(models.Rule)

	ruleError, valid := validators.ValidateRule(rule)

	a.False(valid)
	a.Equal(
		models.ModelValidationError{
			{
				Name: "target.path",
				Errors: []models.ValidationError{
					{TranslationKey: "target.path.invalidPlaceholderIndex"},
				},
			},
		},
		ruleError,
	)
}

func TestValidateRule_InvalidPattern(t *testing.T) {
	a := assert.New(t)

	rule, _ := factories.RuleFactory.MustCreateWithOption(map[string]interface{}{
		"Resolver":   models.RuleResolverPattern,
		"SourcePath": "/r/(\\d+",
	}).(models.Rule)

	ruleError, valid := validators.ValidateRule(rule)

	a.False(valid)
	a.Equal(
		models.ModelValidationError{
			models.FieldValidationError{
				Name: "sourcePath",
				Errors: []models.ValidationError{
					{TranslationKey: "rule.sourcePath.invalidPattern"},
				},
			},
		},
		ruleError,
	)
}
