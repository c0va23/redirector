package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/c0va23/redirector/models"
)

type StoreMock struct {
	mock.Mock
}

func (s *StoreMock) ListHostRules() ([]models.HostRule, error) {
	args := s.MethodCalled("ListHostRules")

	return args.Get(0).([]models.HostRule), args.Error(1)
}

func (s *StoreMock) ReplaceHostRule(hostRule models.HostRule) error {
	args := s.MethodCalled("ReplaceHostRule", hostRule)
	return args.Error(0)
}

func (s *StoreMock) GetHostRules(host string) (*models.HostRule, error) {
	args := s.MethodCalled("GetHostRules", host)

	if err := args.Error(1); nil != err {
		return nil, err
	}

	if hostRule := args.Get(0); nil != hostRule {
		return hostRule.(*models.HostRule), nil
	}
	return nil, nil
}
