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

	r := new(resolver.PatternResolver)
	hostRules := factories.HostRulesFactory.MustCreate().(models.HostRules)
	path := factories.GeneratePath()

	a.Equal(
		hostRules.DefaultTarget,
		r.Resolve(hostRules, path),
	)
}

func TestPatternResolver_MatchSimpleRegexp(t *testing.T) {
	a := assert.New(t)

	r := new(resolver.PatternResolver)
	hostRules := factories.HostRulesFactory.MustCreate().(models.HostRules)

	a.Equal(
		hostRules.Rules[0].Target,
		r.Resolve(hostRules, hostRules.Rules[0].SourcePath),
	)
}

func TestPatternResolver_MatchComplexRegexp(t *testing.T) {
	a := assert.New(t)

	r := new(resolver.PatternResolver)
	hostRules := factories.HostRulesFactory.MustCreateWithOption(map[string]interface{}{
		"Rules": []models.Rule{
			models.Rule{
				SourcePath: "/(one|two)",
				Target: models.Target{
					HTTPCode: 301,
					Path:     "/three",
				},
			},
		},
	}).(models.HostRules)

	a.Equal(
		hostRules.Rules[0].Target,
		r.Resolve(hostRules, "/one"),
	)

	a.Equal(
		hostRules.Rules[0].Target,
		r.Resolve(hostRules, "/two"),
	)

	a.Equal(
		hostRules.DefaultTarget,
		r.Resolve(hostRules, "/other"),
	)
}

func TestPatternResolver_MatchWithPlaceholder(t *testing.T) {
	a := assert.New(t)

	r := new(resolver.PatternResolver)
	var httpCode int32 = 301
	hostRules := factories.HostRulesFactory.MustCreateWithOption(map[string]interface{}{
		"Rules": []models.Rule{
			models.Rule{
				SourcePath: "/u/(\\d+)",
				Target: models.Target{
					HTTPCode: httpCode,
					Path:     "/users/{0}",
				},
			},
		},
	}).(models.HostRules)

	a.Equal(
		models.Target{
			HTTPCode: httpCode,
			Path:     "/users/123",
		},
		r.Resolve(hostRules, "/u/123"),
	)

	a.Equal(
		models.Target{
			HTTPCode: httpCode,
			Path:     "/users/0",
		},
		r.Resolve(hostRules, "/u/0"),
	)

	a.Equal(
		hostRules.DefaultTarget,
		r.Resolve(hostRules, "/u/"),
	)
}
