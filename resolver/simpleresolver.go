package resolver

import (
	"github.com/c0va23/redirector/models"
)

// SimpleResolver is resolver source path by full match (without pattern)
type SimpleResolver struct{}

// NewSimpleResolver create new SimpleResolver
func NewSimpleResolver() SimpleResolver {
	return SimpleResolver{}
}

// Resolve implement Resolver.Resolve
func (r *SimpleResolver) Resolve(
	hostRules models.HostRules,
	sourcePath string,
) models.Target {
	for _, rule := range hostRules.Rules {
		if rule.SourcePath == sourcePath {
			return rule.Target
		}
	}

	return hostRules.DefaultTarget
}
