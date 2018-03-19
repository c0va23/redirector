package controllers_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/c0va23/redirector/controllers"
	"github.com/c0va23/redirector/models"
	"github.com/c0va23/redirector/restapi/operations"

	"github.com/c0va23/redirector/test/factories"
	"github.com/c0va23/redirector/test/mocks"
	"github.com/icrowley/fake"
)

func TestListHostRulesHandler_Success(t *testing.T) {
	a := assert.New(t)

	s := new(mocks.StoreMock)
	r := new(mocks.ResolverMock)
	c := controllers.NewController(s, r)

	listHostRules := make([]models.HostRules, 0, 3)
	for i := 0; i < cap(listHostRules); i++ {
		hostRules := factories.HostRulesFactory.MustCreate().(models.HostRules)
		listHostRules = append(listHostRules, hostRules)
	}
	s.On("ListHostRules").Return(listHostRules, nil)

	a.Equal(
		c.ListHostRulesHandler(operations.ListHostRulesParams{}),
		operations.NewListHostRulesOK().WithPayload(listHostRules),
	)

	s.AssertExpectations(t)
}

func TestListHostRulesHandler_Error(t *testing.T) {
	a := assert.New(t)

	s := new(mocks.StoreMock)
	r := new(mocks.ResolverMock)
	c := controllers.NewController(s, r)

	err := fmt.Errorf("ListHostRulesError")

	s.On("ListHostRules").Return([]models.HostRules{}, err)

	a.Equal(
		c.ListHostRulesHandler(operations.ListHostRulesParams{}),
		operations.NewListHostRulesInternalServerError().
			WithPayload(&models.ServerError{Message: err.Error()}),
	)

	s.AssertExpectations(t)
}

func TestReplaceHostRulesHandler_Success(t *testing.T) {
	a := assert.New(t)

	s := new(mocks.StoreMock)
	r := new(mocks.ResolverMock)
	c := controllers.NewController(s, r)

	newHostRules := factories.HostRulesFactory.MustCreate().(models.HostRules)

	s.On("ReplaceHostRules", newHostRules).Return(nil)

	a.Equal(
		c.ReplaceHostRulesHandler(
			operations.ReplaceHostRulesParams{
				HostRules: newHostRules,
			},
		),
		operations.NewReplaceHostRulesOK().WithPayload(&newHostRules),
	)

	s.AssertExpectations(t)
}

func TestReplaceHostRulesHandler_Error(t *testing.T) {
	a := assert.New(t)

	s := new(mocks.StoreMock)
	r := new(mocks.ResolverMock)
	c := controllers.NewController(s, r)

	newHostRules := factories.HostRulesFactory.MustCreate().(models.HostRules)
	err := fmt.Errorf("ReplaceHostRulesError")
	s.On("ReplaceHostRules", newHostRules).Return(err)

	a.Equal(
		c.ReplaceHostRulesHandler(
			operations.ReplaceHostRulesParams{
				HostRules: newHostRules,
			},
		),
		operations.NewReplaceHostRulesInternalServerError().
			WithPayload(&models.ServerError{Message: err.Error()}),
	)

	s.AssertExpectations(t)
}

func TestRedirectHandler_ServerError(t *testing.T) {
	a := assert.New(t)

	s := new(mocks.StoreMock)
	r := new(mocks.ResolverMock)
	c := controllers.NewController(s, r)

	err := fmt.Errorf("GetHostRulesErr")
	host := fake.DomainName()
	s.On("GetHostRules", host).Return(nil, err)

	a.Equal(
		c.RedirectHandler(operations.RedirectParams{
			Host: host,
		}),
		operations.NewRedirectInternalServerError().
			WithPayload(&models.ServerError{Message: err.Error()}),
	)

	s.AssertExpectations(t)
}

func TestRedirectHandler_NotFound(t *testing.T) {
	a := assert.New(t)

	s := new(mocks.StoreMock)
	r := new(mocks.ResolverMock)
	c := controllers.NewController(s, r)

	host := fake.DomainName()
	s.On("GetHostRules", host).Return(nil, nil)

	a.Equal(
		c.RedirectHandler(operations.RedirectParams{
			Host: host,
		}),
		operations.NewRedirectNotFound(),
	)

	s.AssertExpectations(t)
}

func TestRedirectHandler_Success(t *testing.T) {
	a := assert.New(t)

	s := new(mocks.StoreMock)
	r := new(mocks.ResolverMock)
	c := controllers.NewController(s, r)

	path := factories.GeneratePath()
	hostRules := factories.HostRulesFactory.MustCreate().(models.HostRules)
	s.On("GetHostRules", hostRules.Host).Return(&hostRules, nil)
	r.On("Resolve", hostRules, path).Return(hostRules.DefaultTarget)

	a.Equal(
		c.RedirectHandler(operations.RedirectParams{
			Host:       hostRules.Host,
			SourcePath: path,
		}),
		operations.NewRedirectDefault(int(hostRules.DefaultTarget.HTTPCode)).
			WithLocation(hostRules.DefaultTarget.Path),
	)

	s.AssertExpectations(t)
}
