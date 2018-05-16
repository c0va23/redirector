package controllers

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/c0va23/redirector/models"
	"github.com/c0va23/redirector/resolver"
	"github.com/c0va23/redirector/restapi/operations/config"
	"github.com/c0va23/redirector/restapi/operations/redirect"
	"github.com/c0va23/redirector/store"
)

// Controller implement methods into restapi
type Controller struct {
	store    store.Store
	resolver resolver.Resolver
}

// NewController initialize new controller
func NewController(store store.Store, resolver resolver.Resolver) Controller {
	return Controller{
		store:    store,
		resolver: resolver,
	}
}

// ListHostRulesHandler is handler for ListHostRules
func (c *Controller) ListHostRulesHandler(params config.ListHostRulesParams, _principal interface{}) middleware.Responder {
	listHostRules, err := c.store.ListHostRules()

	if nil != err {
		serverError := models.ServerError{Message: err.Error()}
		return config.NewListHostRulesInternalServerError().WithPayload(&serverError)
	}

	return config.NewListHostRulesOK().
		WithPayload(listHostRules)
}

// CreateHostRulesHandler is handler for CreateHostRules
func (c *Controller) CreateHostRulesHandler(
	params config.CreateHostRulesParams,
	_principal interface{},
) middleware.Responder {
	err := c.store.CreateHostRules(params.HostRules)

	if nil != err {
		switch err {
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

	return config.
		NewCreateHostRulesOK().
		WithPayload(&params.HostRules)
}

// UpdateHostRulesHandler is handler for UpdateHostRules
func (c *Controller) UpdateHostRulesHandler(
	params config.UpdateHostRulesParams,
	_principal interface{},
) middleware.Responder {
	err := c.store.UpdateHostRules(params.Host, params.HostRules)

	if nil != err {
		switch err {
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

	return config.
		NewUpdateHostRulesOK().
		WithPayload(&params.HostRules)
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

// RedirectHandler is handler for Redirect
func (c *Controller) RedirectHandler(params redirect.RedirectParams) middleware.Responder {
	hostRules, err := c.store.GetHostRules(params.Host)

	switch err {
	case nil:
		target := c.resolver.Resolve(*hostRules, params.SourcePath)
		return redirect.NewRedirectDefault(int(target.HTTPCode)).
			WithLocation(target.Path)
	case store.ErrNotFound:
		return redirect.NewRedirectNotFound()
	default:
		return redirect.NewRedirectInternalServerError().
			WithPayload(&models.ServerError{Message: err.Error()})
	}

}
