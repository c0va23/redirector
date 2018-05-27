package mocks_test

import (
	"testing"

	"github.com/c0va23/redirector/store"
	"github.com/c0va23/redirector/test/mocks"

	"github.com/stretchr/testify/assert"
)

func TestStoreMock(t *testing.T) {
	assert.Implements(t, (*store.Store)(nil), &mocks.StoreMock{})
}
