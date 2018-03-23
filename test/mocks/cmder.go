package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/mediocregopher/radix.v2/redis"
)

// CmderMock is mock for Redis client
type CmderMock struct {
	mock.Mock
}

// Cmd implement Cmder.Cmd
func (c *CmderMock) Cmd(cmd string, args ...interface{}) *redis.Resp {
	return c.Mock.MethodCalled("Cmd", cmd, args)[0].(*redis.Resp)
}
