package mocks

import (
	"github.com/c0va23/redirector/models"
	"github.com/stretchr/testify/mock"
)

// ResolverMock is mock for resolver.Resolver
type ResolverMock struct {
	mock.Mock
}

// Resolve implement rosolver.Resolver.Resolve
func (r *ResolverMock) Resolve(hostRules models.HostRules, sourcePath string) models.Target {
	args := r.MethodCalled("Resolve", hostRules, sourcePath)
	return args.Get(0).(models.Target)
}
