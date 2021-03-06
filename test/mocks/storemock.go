package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/c0va23/redirector/models"
)

// StoreMock is mock for Store interface
type StoreMock struct {
	mock.Mock
}

// ListHostRules for store mock
func (s *StoreMock) ListHostRules() ([]models.HostRules, error) {
	args := s.MethodCalled("ListHostRules")

	return args.Get(0).([]models.HostRules), args.Error(1)
}

// CreateHostRules for store mock
func (s *StoreMock) CreateHostRules(hostRules models.HostRules) error {
	args := s.MethodCalled("CreateHostRules", hostRules)
	return args.Error(0)
}

// UpdateHostRules for store mock
func (s *StoreMock) UpdateHostRules(host string, hostRules models.HostRules) error {
	args := s.MethodCalled("UpdateHostRules", host, hostRules)
	return args.Error(0)
}

// DeleteHostRules for store mock
func (s *StoreMock) DeleteHostRules(host string) error {
	args := s.MethodCalled("DeleteHostRules", host)
	return args.Error(0)
}

// GetHostRules for store mock
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

// CheckHealth for store mock
func (s *StoreMock) CheckHealth() error {
	args := s.MethodCalled("CheckHealth")

	return args.Error(0)
}
