package controllers

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/c0va23/redirector/models"
	"github.com/c0va23/redirector/restapi/operations/config"
	"github.com/c0va23/redirector/restapi/operations/redirect"
	"github.com/c0va23/redirector/store"
)

// Controller implement methods into restapi
type Controller struct {
	store store.Store
}

// NewController initialize new controller
func NewController(store store.Store) Controller {
	return Controller{
		store: store,
	}
}

// ListHostRulesHandler is handler for ListHostRules
func (c *Controller) ListHostRulesHandler(params config.ListHostRulesParams, _principal interface{}) middleware.Responder {
	listHostRules, err := c.store.ListHostRules()

	switch err {
	case nil:
		return config.NewListHostRulesOK().
			WithPayload(listHostRules)
	default:
		return config.
			NewListHostRulesInternalServerError().
			WithPayload(&models.ServerError{Message: err.Error()})
	}
}

// CreateHostRulesHandler is handler for CreateHostRules
func (c *Controller) CreateHostRulesHandler(
	params config.CreateHostRulesParams,
	_principal interface{},
) middleware.Responder {
	err := c.store.CreateHostRules(params.HostRules)

	switch err {
	case nil:
		return config.
			NewCreateHostRulesOK().
			WithPayload(&params.HostRules)
	case store.ErrExists:
		return config.NewCreateHostRulesConflict()
	default:
		return config.
			NewCreateHostRulesInternalServerError().
			WithPayload(&models.ServerError{
				Message: err.Error(),
			})

	}
}

// UpdateHostRulesHandler is handler for UpdateHostRules
func (c *Controller) UpdateHostRulesHandler(
	params config.UpdateHostRulesParams,
	_principal interface{},
) middleware.Responder {
	err := c.store.UpdateHostRules(params.Host, params.HostRules)

	switch err {
	case nil:
		return config.
			NewUpdateHostRulesOK().
			WithPayload(&params.HostRules)
	case store.ErrNotFound:
		return config.NewUpdateHostRulesNotFound()
	default:
		return config.
			NewUpdateHostRulesInternalServerError().
			WithPayload(&models.ServerError{
				Message: err.Error(),
			})
	}
}

// GetHostRulesHandler is handler for GetHostRules
func (c *Controller) GetHostRulesHandler(
	params config.GetHostRuleParams,
	_principal interface{},
) middleware.Responder {
	hostRules, err := c.store.GetHostRules(params.Host)

	switch err {
	case nil:
		return config.NewGetHostRuleOK().
			WithPayload(hostRules)
	case store.ErrNotFound:
		return config.NewGetHostRuleNotFound()
	default:
		return config.NewGetHostRuleInternalServerError().
			WithPayload(&models.ServerError{Message: err.Error()})
	}
}

// DeleteHostRulesHandler is handler for DeleteHostRules
func (c *Controller) DeleteHostRulesHandler(
	params config.DeleteHostRulesParams,
	_principal interface{},
) middleware.Responder {
	err := c.store.DeleteHostRules(params.Host)
	switch err {
	case nil:
		return config.NewDeleteHostRulesNoContent()
	case store.ErrNotFound:
		return config.NewDeleteHostRulesNotFound()
	default:
		return config.
			NewDeleteHostRulesInternalServerError().
			WithPayload(&models.ServerError{
				Message: err.Error(),
			})

	}
}

// HealthCheckHandler is handler for HealthCheck
func (c *Controller) HealthCheckHandler(params redirect.HealthcheckParams) middleware.Responder {
	err := c.store.CheckHealth()

	switch err {
	case nil:
		return redirect.NewHealthcheckOK()
	default:
		return redirect.NewHealthcheckInternalServerError()
	}
}
