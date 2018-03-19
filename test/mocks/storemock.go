package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/c0va23/redirector/models"
)

type StoreMock struct {
	mock.Mock
}

func (s *StoreMock) ListHostRules() ([]models.HostRules, error) {
	args := s.MethodCalled("ListHostRules")

	return args.Get(0).([]models.HostRules), args.Error(1)
}

func (s *StoreMock) ReplaceHostRules(hostRules models.HostRules) error {
	args := s.MethodCalled("ReplaceHostRules", hostRules)
	return args.Error(0)
}

func (s *StoreMock) GetHostRules(host string) (*models.HostRules, error) {
	args := s.MethodCalled("GetHostRules", host)

	if err := args.Error(1); nil != err {
		return nil, err
	}

	if hostRules := args.Get(0); nil != hostRules {
		return hostRules.(*models.HostRules), nil
	}
	return nil, nil
}
