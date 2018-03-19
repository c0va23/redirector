package resolver_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/c0va23/redirector/models"
	"github.com/c0va23/redirector/resolver"

	"github.com/c0va23/redirector/test/factories"
)

func TestSimpleResolver_NotMatchPath(t *testing.T) {
	a := assert.New(t)

	r := resolver.NewSimpleResolver()
	hostRules := factories.HostRulesFactory.MustCreate().(models.HostRules)
	path := factories.GeneratePath()

	a.Equal(
		r.Resolve(hostRules, path),
		hostRules.DefaultTarget,
	)
}

func TestSimpleResolver_MatchPath(t *testing.T) {
	a := assert.New(t)

	r := resolver.NewSimpleResolver()
	hostRules := factories.HostRulesFactory.MustCreate().(models.HostRules)

	rule := hostRules.Rules[0]

	a.Equal(
		r.Resolve(hostRules, rule.SourcePath),
		rule.Target,
	)
}
