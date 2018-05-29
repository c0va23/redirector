package resolver_test

import (
	"testing"

	"github.com/c0va23/redirector/models"
	"github.com/c0va23/redirector/resolver"
	"github.com/c0va23/redirector/test/factories"

	"github.com/stretchr/testify/assert"
)

func TestResolve_WithoutResolvers(t *testing.T) {
	a := assert.New(t)

	hostRules := factories.HostRulesFactory.MustCreate().(models.HostRules)
	emptyResolvers := resolver.MultiHostRulesResolver{}

	a.Equal(
		hostRules.DefaultTarget,
		emptyResolvers.Resolve(hostRules, "/"),
	)
}

func TestResolve_ResolveNotMatch(t *testing.T) {
	a := assert.New(t)

	rule := factories.RuleFactory.MustCreate().(models.Rule)
	hostRules := factories.HostRulesFactory.
		MustCreateWithOption(map[string]interface{}{
			"Rules": []models.Rule{rule},
		}).(models.HostRules)

	ruleResolver := func(rule models.Rule, sourcePath string) *models.Target {
		return nil
	}

	hostRulesResolver := resolver.MultiHostRulesResolver{
		"fakeResolver": ruleResolver,
	}

	a.Equal(
		hostRules.DefaultTarget,
		hostRulesResolver.Resolve(hostRules, "/"),
	)
}

func TestResolve_ResolveMatch(t *testing.T) {
	a := assert.New(t)

	rule := factories.RuleFactory.MustCreate().(models.Rule)
	hostRules := factories.HostRulesFactory.
		MustCreateWithOption(map[string]interface{}{
			"Rules": []models.Rule{rule},
		}).(models.HostRules)

	ruleResolver := func(rule models.Rule, sourcePath string) *models.Target {
		return &rule.Target
	}

	hostRulesResolver := resolver.MultiHostRulesResolver{
		rule.Resolver: ruleResolver,
	}

	a.Equal(
		rule.Target,
		hostRulesResolver.Resolve(hostRules, "/"),
	)
}
