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
func (r *ResolverMock) Resolve(hostRule models.HostRule, sourcePath string) models.Target {
	args := r.MethodCalled("Resolve", hostRule, sourcePath)
	return args.Get(0).(models.Target)
}
