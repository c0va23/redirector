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
	c := controllers.NewController(s)

	hostRules := make([]models.HostRule, 0, 3)
	for i := 0; i < cap(hostRules); i++ {
		hostRule := factories.HostRuleFactory.MustCreate().(models.HostRule)
		hostRules = append(hostRules, hostRule)
	}
	s.On("ListHostRules").Return(hostRules, nil)

	a.Equal(
		c.ListHostRulesHandler(operations.ListHostRulesParams{}),
		operations.NewListHostRulesOK().WithPayload(hostRules),
	)

	s.AssertExpectations(t)
}

func TestListHostRulesHandler_Error(t *testing.T) {
	a := assert.New(t)

	s := new(mocks.StoreMock)
	c := controllers.NewController(s)

	err := fmt.Errorf("ListHostRulesError")

	s.On("ListHostRules").Return([]models.HostRule{}, err)

	a.Equal(
		c.ListHostRulesHandler(operations.ListHostRulesParams{}),
		operations.NewListHostRulesInternalServerError().
			WithPayload(&models.ServerError{Message: err.Error()}),
	)

	s.AssertExpectations(t)
}

func TestReplaceHostRuleHandler_Success(t *testing.T) {
	a := assert.New(t)

	s := new(mocks.StoreMock)
	c := controllers.NewController(s)

	newHostRule := factories.HostRuleFactory.MustCreate().(models.HostRule)

	s.On("ReplaceHostRule", newHostRule).Return(nil)

	a.Equal(
		c.ReplaceHostRulesHandler(
			operations.ReplaceHostRuleParams{
				HostRule: newHostRule,
			},
		),
		operations.NewReplaceHostRuleOK().WithPayload(&newHostRule),
	)

	s.AssertExpectations(t)
}

func TestReplaceHostRuleHandler_Error(t *testing.T) {
	a := assert.New(t)

	s := new(mocks.StoreMock)
	c := controllers.NewController(s)

	newHostRule := factories.HostRuleFactory.MustCreate().(models.HostRule)
	err := fmt.Errorf("ReplaceHostRuleError")
	s.On("ReplaceHostRule", newHostRule).Return(err)

	a.Equal(
		c.ReplaceHostRulesHandler(
			operations.ReplaceHostRuleParams{
				HostRule: newHostRule,
			},
		),
		operations.NewReplaceHostRuleInternalServerError().
			WithPayload(&models.ServerError{Message: err.Error()}),
	)

	s.AssertExpectations(t)
}

func TestRedirectHandler_ServerError(t *testing.T) {
	a := assert.New(t)

	s := new(mocks.StoreMock)
	c := controllers.NewController(s)

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
	c := controllers.NewController(s)

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
	c := controllers.NewController(s)

	path := factories.GeneratePath()
	hostRule := factories.HostRuleFactory.MustCreate().(models.HostRule)
	s.On("GetHostRules", hostRule.Host).Return(&hostRule, nil)

	a.Equal(
		c.RedirectHandler(operations.RedirectParams{
			Host:       hostRule.Host,
			SourcePath: path,
		}),
		operations.NewRedirectDefault(int(hostRule.DefaultTarget.HTTPCode)).
			WithLocation(hostRule.DefaultTarget.TargetPath),
	)

	s.AssertExpectations(t)
}
