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

	rule := factories.RuleFactory.MustCreate().(models.Rule)
	path := factories.GeneratePath()

	a.Nil(resolver.SimpleResolver(rule, path))
}

func TestSimpleResolver_MatchPath(t *testing.T) {
	a := assert.New(t)

	rule := factories.RuleFactory.MustCreate().(models.Rule)

	a.Equal(
		&rule.Target,
		resolver.SimpleResolver(rule, rule.SourcePath),
	)
}
