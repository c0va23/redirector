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
	hostRule := factories.HostRuleFactory.MustCreate().(models.HostRule)
	path := factories.GeneratePath()

	a.Equal(
		r.Resolve(hostRule, path),
		hostRule.DefaultTarget,
	)
}

func TestSimpleResolver_MatchPath(t *testing.T) {
	a := assert.New(t)

	r := resolver.NewSimpleResolver()
	hostRule := factories.HostRuleFactory.MustCreate().(models.HostRule)

	rule := hostRule.Rules[0]

	a.Equal(
		r.Resolve(hostRule, rule.SourcePathPattern),
		rule.Target,
	)
}
