package mocks_test

import (
	"testing"

	"github.com/mediocregopher/radix.v2/util"

	"github.com/stretchr/testify/assert"

	"github.com/c0va23/redirector/test/mocks"
)

func TestCmderMock(t *testing.T) {
	a := assert.New(t)

	a.Implements((*util.Cmder)(nil), new(mocks.CmderMock))
}
