package handlers_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/c0va23/redirector/handlers"
	"github.com/c0va23/redirector/models"
	"github.com/c0va23/redirector/restapi/operations/config"
	"github.com/c0va23/redirector/restapi/operations/redirect"
	"github.com/c0va23/redirector/validators"

	"github.com/c0va23/redirector/store"
	"github.com/c0va23/redirector/test/factories"
	"github.com/c0va23/redirector/test/mocks"
	"github.com/icrowley/fake"
)

func TestListHostRulesHandler_Success(t *testing.T) {
	a := assert.New(t)

	s := new(mocks.StoreMock)
	ch := handlers.NewConfigHandlers(s)

	listHostRules := make([]models.HostRules, 0, 3)
	for i := 0; i < cap(listHostRules); i++ {
		hostRules := factories.
			HostRulesFactory.
			MustCreateWithOption(map[string]interface{}{
				"Host": fmt.Sprintf("host%d.org", cap(listHostRules)-i),
			}).(models.HostRules)
		listHostRules = append(listHostRules, hostRules)
	}
	s.On("ListHostRules").Return(listHostRules, nil)

	revertHostRules := make([]models.HostRules, 0, 3)
	for i := len(listHostRules) - 1; i >= 0; i-- {
		revertHostRules = append(revertHostRules, listHostRules[i])
	}

	a.Equal(
		ch.ListHostRulesHandler(config.ListHostRulesParams{}, true),
		config.NewListHostRulesOK().WithPayload(revertHostRules),
	)

	s.AssertExpectations(t)
}

func TestListHostRulesHandler_Error(t *testing.T) {
	a := assert.New(t)

	s := new(mocks.StoreMock)
	ch := handlers.NewConfigHandlers(s)

	err := fmt.Errorf("ListHostRulesError")

	s.On("ListHostRules").Return([]models.HostRules{}, err)

	a.Equal(
		ch.ListHostRulesHandler(config.ListHostRulesParams{}, true),
		config.NewListHostRulesInternalServerError().
			WithPayload(&models.ServerError{Message: err.Error()}),
	)

	s.AssertExpectations(t)
}

func TestCreateHostRulesHandler_Success(t *testing.T) {
	a := assert.New(t)

	s := new(mocks.StoreMock)
	ch := handlers.NewConfigHandlers(s)

	newHostRules := factories.HostRulesFactory.MustCreate().(models.HostRules)

	s.On("CreateHostRules", newHostRules).Return(nil)

	a.Equal(
		config.NewCreateHostRulesOK().
			WithPayload(&newHostRules),
		ch.CreateHostRulesHandler(
			config.CreateHostRulesParams{
				HostRules: newHostRules,
			},
			true,
		),
	)

	s.AssertExpectations(t)
}

func TestCreateHostRulesHandler_ExistsError(t *testing.T) {
	a := assert.New(t)

	s := new(mocks.StoreMock)
	ch := handlers.NewConfigHandlers(s)

	newHostRules := factories.HostRulesFactory.MustCreate().(models.HostRules)
	s.On("CreateHostRules", newHostRules).Return(store.ErrExists)

	a.Equal(
		config.NewCreateHostRulesConflict(),
		ch.CreateHostRulesHandler(
			config.CreateHostRulesParams{
				HostRules: newHostRules,
			},
			true,
		),
	)

	s.AssertExpectations(t)
}

func TestCreateHostRulesHandler_OtherError(t *testing.T) {
	a := assert.New(t)

	s := new(mocks.StoreMock)
	ch := handlers.NewConfigHandlers(s)

	newHostRules := factories.HostRulesFactory.MustCreate().(models.HostRules)
	err := fmt.Errorf("CreateHostRulesError")
	s.On("CreateHostRules", newHostRules).Return(err)

	a.Equal(
		config.NewCreateHostRulesInternalServerError().
			WithPayload(&models.ServerError{Message: err.Error()}),
		ch.CreateHostRulesHandler(
			config.CreateHostRulesParams{
				HostRules: newHostRules,
			},
			true,
		),
	)

	s.AssertExpectations(t)
}

func TestCreateHostRulesHandler_ValidationError(t *testing.T) {
	a := assert.New(t)

	s := new(mocks.StoreMock)
	ch := handlers.NewConfigHandlers(s)

	newHostRules := models.HostRules{}
	hostRulesError, _ := validators.ValidateHostRules(newHostRules)

	a.Equal(
		config.
			NewCreateHostRulesUnprocessableEntity().
			WithPayload(hostRulesError),
		ch.CreateHostRulesHandler(
			config.CreateHostRulesParams{
				HostRules: newHostRules,
			},
			true,
		),
	)

	s.AssertExpectations(t)
}

func TestUpdateHostRules_Success(t *testing.T) {
	a := assert.New(t)

	s := new(mocks.StoreMock)
	ch := handlers.NewConfigHandlers(s)

	host := fake.DomainName()
	hostRules := factories.HostRulesFactory.MustCreate().(models.HostRules)

	s.On("UpdateHostRules", host, hostRules).Return(nil)

	a.Equal(
		config.NewUpdateHostRulesOK().
			WithPayload(&hostRules),
		ch.UpdateHostRulesHandler(
			config.UpdateHostRulesParams{
				Host:      host,
				HostRules: hostRules,
			},
			true,
		),
	)

	s.AssertExpectations(t)
}

func TestUpdateHostRules_NotFoundError(t *testing.T) {
	a := assert.New(t)

	s := new(mocks.StoreMock)
	ch := handlers.NewConfigHandlers(s)

	host := fake.DomainName()
	hostRules := factories.HostRulesFactory.MustCreate().(models.HostRules)

	s.On("UpdateHostRules", host, hostRules).Return(store.ErrNotFound)

	a.Equal(
		config.NewUpdateHostRulesNotFound(),
		ch.UpdateHostRulesHandler(
			config.UpdateHostRulesParams{
				Host:      host,
				HostRules: hostRules,
			},
			true,
		),
	)

	s.AssertExpectations(t)
}

func TestUpdateHostRules_OtherError(t *testing.T) {
	a := assert.New(t)

	s := new(mocks.StoreMock)
	ch := handlers.NewConfigHandlers(s)

	host := fake.DomainName()
	hostRules := factories.HostRulesFactory.MustCreate().(models.HostRules)
	err := fmt.Errorf("UpdateHostRulesError")

	s.On("UpdateHostRules", host, hostRules).Return(err)

	a.Equal(
		config.NewUpdateHostRulesInternalServerError().
			WithPayload(&models.ServerError{
				Message: err.Error(),
			}),
		ch.UpdateHostRulesHandler(
			config.UpdateHostRulesParams{
				Host:      host,
				HostRules: hostRules,
			},
			true,
		),
	)

	s.AssertExpectations(t)
}

func TestUpdateHostRules_ValidationError(t *testing.T) {
	a := assert.New(t)

	s := new(mocks.StoreMock)
	ch := handlers.NewConfigHandlers(s)

	host := fake.DomainName()
	hostRules := models.HostRules{}

	hostRulesError, _ := validators.ValidateHostRules(hostRules)

	a.Equal(
		config.NewUpdateHostRulesUnprocessableEntity().
			WithPayload(hostRulesError),
		ch.UpdateHostRulesHandler(
			config.UpdateHostRulesParams{
				Host:      host,
				HostRules: hostRules,
			},
			true,
		),
	)

	s.AssertExpectations(t)
}

func TestGetHostRules_Error(t *testing.T) {
	a := assert.New(t)

	host := fake.DomainName()
	err := fmt.Errorf("GetHostRulesErr")

	s := new(mocks.StoreMock)
	s.On("GetHostRules", host).Return(nil, err)

	ch := handlers.NewConfigHandlers(s)

	a.Equal(
		ch.GetHostRulesHandler(
			config.GetHostRuleParams{
				Host: host,
			},
			true,
		),
		config.NewGetHostRuleInternalServerError().
			WithPayload(&models.ServerError{Message: err.Error()}),
	)

	s.AssertExpectations(t)
}

func TestGetHostRules_NotFound(t *testing.T) {
	a := assert.New(t)

	host := fake.DomainName()
	s := new(mocks.StoreMock)
	s.On("GetHostRules", host).Return(nil, store.ErrNotFound)

	ch := handlers.NewConfigHandlers(s)

	a.Equal(
		ch.GetHostRulesHandler(
			config.GetHostRuleParams{
				Host: host,
			},
			true,
		),
		config.NewGetHostRuleNotFound(),
	)

	s.AssertExpectations(t)
}

func TestGetHostRules_Found(t *testing.T) {
	a := assert.New(t)

	hostRules := factories.
		HostRulesFactory.
		MustCreate().(models.HostRules)

	s := new(mocks.StoreMock)
	s.On("GetHostRules", hostRules.Host).Return(&hostRules, nil)

	ch := handlers.NewConfigHandlers(s)

	a.Equal(
		ch.GetHostRulesHandler(
			config.GetHostRuleParams{
				Host: hostRules.Host,
			},
			true,
		),
		config.NewGetHostRuleOK().
			WithPayload(&hostRules),
	)

	s.AssertExpectations(t)
}

func TestDeleteHostRules_Success(t *testing.T) {
	a := assert.New(t)

	host := fake.DomainName()

	s := new(mocks.StoreMock)
	s.On("DeleteHostRules", host).Return(nil)

	ch := handlers.NewConfigHandlers(s)

	a.Equal(
		ch.DeleteHostRulesHandler(
			config.DeleteHostRulesParams{
				Host: host,
			},
			true,
		),
		config.NewDeleteHostRulesNoContent(),
	)

	s.AssertExpectations(t)
}

func TestDeleteHostRules_NotFoundError(t *testing.T) {
	a := assert.New(t)

	host := fake.DomainName()

	s := new(mocks.StoreMock)
	s.On("DeleteHostRules", host).Return(store.ErrNotFound)

	ch := handlers.NewConfigHandlers(s)

	a.Equal(
		ch.DeleteHostRulesHandler(
			config.DeleteHostRulesParams{
				Host: host,
			},
			true,
		),
		config.NewDeleteHostRulesNotFound(),
	)

	s.AssertExpectations(t)
}

func TestDeleteHostRules_OtherError(t *testing.T) {
	a := assert.New(t)

	host := fake.DomainName()

	otherErr := fmt.Errorf("DeleteErr")
	s := new(mocks.StoreMock)
	s.On("DeleteHostRules", host).Return(otherErr)

	ch := handlers.NewConfigHandlers(s)

	a.Equal(
		ch.DeleteHostRulesHandler(
			config.DeleteHostRulesParams{
				Host: host,
			},
			true,
		),
		config.NewDeleteHostRulesInternalServerError().
			WithPayload(&models.ServerError{
				Message: otherErr.Error(),
			}),
	)

	s.AssertExpectations(t)
}

func TestHealthCheckHandler_Success(t *testing.T) {
	a := assert.New(t)

	s := new(mocks.StoreMock)

	ch := handlers.NewConfigHandlers(s)

	s.On("CheckHealth").Return(nil)

	a.Equal(
		redirect.NewHealthcheckOK(),
		ch.HealthCheckHandler(redirect.NewHealthcheckParams()),
	)
}

func TestHealthCheckHandler_Error(t *testing.T) {
	a := assert.New(t)

	s := new(mocks.StoreMock)

	ch := handlers.NewConfigHandlers(s)

	err := fmt.Errorf("CheckHealthError")
	s.On("CheckHealth").Return(err)

	a.Equal(
		redirect.NewHealthcheckInternalServerError(),
		ch.HealthCheckHandler(redirect.NewHealthcheckParams()),
	)
}
