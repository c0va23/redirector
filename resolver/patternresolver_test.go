package resolver_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/c0va23/redirector/models"
	"github.com/c0va23/redirector/resolver"

	"github.com/c0va23/redirector/test/factories"
)

func TestPatternResolver_NotMatchPath(t *testing.T) {
	a := assert.New(t)

	rule := factories.RuleFactory.MustCreate().(models.Rule)
	path := factories.GeneratePath()

	a.Nil(resolver.PatternResolver(rule, path))
}

func TestPatternResolver_MatchSimpleRegexp(t *testing.T) {
	a := assert.New(t)

	rule := factories.RuleFactory.MustCreate().(models.Rule)

	a.Equal(
		&rule.Target,
		resolver.PatternResolver(rule, rule.SourcePath),
	)
}

func TestPatternResolver_MatchComplexRegexp(t *testing.T) {
	a := assert.New(t)

	rule := factories.RuleFactory.MustCreateWithOption(map[string]interface{}{
		"SourcePath": "/(one|two)",
		"Target": models.Target{
			HTTPCode: 301,
			Path:     "/three",
		},
	}).(models.Rule)

	a.Equal(
		&rule.Target,
		resolver.PatternResolver(rule, "/one"),
	)

	a.Equal(
		&rule.Target,
		resolver.PatternResolver(rule, "/two"),
	)

	a.Nil(resolver.PatternResolver(rule, "/other"))
}

func TestPatternResolver_MatchWithPlaceholder(t *testing.T) {
	a := assert.New(t)

	var httpCode int32 = 301
	rule := factories.RuleFactory.MustCreateWithOption(map[string]interface{}{
		"SourcePath": "/u/(\\d+)",
		"Target": models.Target{
			HTTPCode: httpCode,
			Path:     "/users/{0}",
		},
	}).(models.Rule)

	a.Equal(
		&models.Target{
			HTTPCode: httpCode,
			Path:     "/users/123",
		},
		resolver.PatternResolver(rule, "/u/123"),
	)

	a.Equal(
		&models.Target{
			HTTPCode: httpCode,
			Path:     "/users/0",
		},
		resolver.PatternResolver(rule, "/u/0"),
	)

	a.Nil(resolver.PatternResolver(rule, "/u/"))
}
