package resolvers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/c0va23/redirector/models"
	"github.com/c0va23/redirector/resolvers"

	"github.com/c0va23/redirector/test/factories"
)

func TestSimpleResolver_NotMatchPath(t *testing.T) {
	a := assert.New(t)

	rule := factories.RuleFactory.MustCreate().(models.Rule)
	path := factories.GeneratePath()

	a.Nil(resolvers.SimpleResolver(rule, path))
}

func TestSimpleResolver_MatchPath(t *testing.T) {
	a := assert.New(t)

	rule := factories.RuleFactory.MustCreate().(models.Rule)

	a.Equal(
		&rule.Target,
		resolvers.SimpleResolver(rule, rule.SourcePath),
	)
}
