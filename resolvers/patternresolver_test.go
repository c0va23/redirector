package resolvers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/c0va23/redirector/models"
	"github.com/c0va23/redirector/resolvers"

	"github.com/c0va23/redirector/test/factories"
)

func TestPatternResolver_NotMatchPath(t *testing.T) {
	a := assert.New(t)

	rule := factories.RuleFactory.MustCreate().(models.Rule)
	path := factories.GeneratePath()

	a.Nil(resolvers.PatternResolver(rule, path))
}

func TestPatternResolver_MatchSimpleRegexp(t *testing.T) {
	a := assert.New(t)

	rule := factories.RuleFactory.MustCreate().(models.Rule)

	a.Equal(
		&rule.Target,
		resolvers.PatternResolver(rule, rule.SourcePath),
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
		resolvers.PatternResolver(rule, "/one"),
	)

	a.Equal(
		&rule.Target,
		resolvers.PatternResolver(rule, "/two"),
	)

	a.Nil(resolvers.PatternResolver(rule, "/other"))
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
		resolvers.PatternResolver(rule, "/u/123"),
	)

	a.Equal(
		&models.Target{
			HTTPCode: httpCode,
			Path:     "/users/0",
		},
		resolvers.PatternResolver(rule, "/u/0"),
	)

	a.Nil(resolvers.PatternResolver(rule, "/u/"))
}

func TestPatternResolver_InvalidRegexp(t *testing.T) {
	a := assert.New(t)

	rule := factories.RuleFactory.MustCreateWithOption(map[string]interface{}{
		"SourcePath": "\\",
	}).(models.Rule)

	a.Nil(resolvers.PatternResolver(rule, "/test"))
}

func TestPatternResolver_MultipleMatch(t *testing.T) {
	a := assert.New(t)

	rule := factories.RuleFactory.MustCreateWithOption(map[string]interface{}{
		"SourcePath": ".",
	}).(models.Rule)

	a.Nil(resolvers.PatternResolver(rule, "/test"))
}
