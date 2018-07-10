package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"

	"github.com/c0va23/redirector/handlers"
	"github.com/c0va23/redirector/models"
	"github.com/c0va23/redirector/store"
	"github.com/c0va23/redirector/test/factories"
	"github.com/c0va23/redirector/test/mocks"
)

func TestNewRedirectHandler(t *testing.T) {
	s := new(mocks.StoreMock)
	r := new(mocks.ResolverMock)
	h := handlers.NewRedirectHandler(s, r)

	assert.Implements(t, (*http.Handler)(nil), h)
}

func TestRedirectHandler_Handle_ServerError(t *testing.T) {
	a := assert.New(t)

	s := new(mocks.StoreMock)
	r := new(mocks.ResolverMock)
	h := handlers.NewRedirectHandler(s, r)

	err := fmt.Errorf("GetHostRulesErr")
	host := fake.DomainName()
	s.On("GetHostRules", host).Return(nil, err)

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Host = host

	h.ServeHTTP(rw, req)

	res := rw.Result()

	a.Equal(http.StatusInternalServerError, res.StatusCode)

	s.AssertExpectations(t)
}

func TestRedirectHandler_Handle_NotFound(t *testing.T) {
	a := assert.New(t)

	s := new(mocks.StoreMock)
	r := new(mocks.ResolverMock)
	h := handlers.NewRedirectHandler(s, r)

	host := fake.DomainName()
	s.On("GetHostRules", host).Return(nil, store.ErrNotFound)

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Host = host

	h.ServeHTTP(rw, req)

	res := rw.Result()

	a.Equal(http.StatusNotFound, res.StatusCode)

	s.AssertExpectations(t)
}

func TestRedirectHandler_Handle_Success(t *testing.T) {
	a := assert.New(t)

	s := new(mocks.StoreMock)
	r := new(mocks.ResolverMock)
	h := handlers.NewRedirectHandler(s, r)

	path := factories.GeneratePath()
	hostRules := factories.HostRulesFactory.MustCreate().(models.HostRules)
	s.On("GetHostRules", hostRules.Host).Return(&hostRules, nil)
	r.On("Resolve", hostRules, path).Return(hostRules.DefaultTarget)

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", path, nil)
	req.Host = hostRules.Host

	h.ServeHTTP(rw, req)

	res := rw.Result()

	a.Equal(int(hostRules.DefaultTarget.HTTPCode), res.StatusCode)
	a.Equal(hostRules.DefaultTarget.Path, res.Header.Get("Location"))

	s.AssertExpectations(t)
}
