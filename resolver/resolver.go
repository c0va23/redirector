package resolver

import (
	"github.com/c0va23/redirector/models"
)

// RuleResolver used for resolve host and path to HTTP code and new URL
type RuleResolver func(models.Rule, string) *models.Target

// HostRulesResolver is interface host HostRules resolver
type HostRulesResolver interface {
	Resolve(models.HostRules, string) models.Target
}

// DefaultResolvers is default resolvers for MultiHostRulesResolver
var DefaultResolvers = map[string]RuleResolver{
	"simple":  SimpleResolver,
	"pattern": PatternResolver,
}

// MultiHostRulesResolver is wrapper for multpe resovler
type MultiHostRulesResolver map[string]RuleResolver

// Resolve match each rule by resolver and reutrn if it return target.
// Otherwise return default target.
func (r MultiHostRulesResolver) Resolve(
	hostRules models.HostRules,
	sourcePath string,
) models.Target {
	for _, rule := range hostRules.Rules {
		for resolverType, resolver := range r {
			if resolverType != rule.Resolver {
				continue
			}

			if target := resolver(rule, sourcePath); nil != target {
				return *target
			}
		}
	}
	return hostRules.DefaultTarget
}
