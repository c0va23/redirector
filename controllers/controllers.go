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

// ReplaceHostRulesHandler is handler for ReplaceHostRules
func (c *Controller) ReplaceHostRulesHandler(params config.ReplaceHostRulesParams, _principal interface{}) middleware.Responder {
	err := c.store.ReplaceHostRules(params.HostRules)

	if nil != err {
		serverError := models.ServerError{Message: err.Error()}
		return config.NewReplaceHostRulesInternalServerError().WithPayload(&serverError)
	}

	return config.NewReplaceHostRulesOK().
		WithPayload(&params.HostRules)
}

// GetHostRulesHandler is handler for GetHostRules
func (c *Controller) GetHostRulesHandler(
	params config.GetHostRuleParams,
	_principal interface{},
) middleware.Responder {
	hostRules, err := c.store.GetHostRules(params.Host)

	if nil != err {
		return config.NewGetHostRuleInternalServerError().
			WithPayload(&models.ServerError{Message: err.Error()})
	}

	if nil == hostRules {
		return config.NewGetHostRuleNotFound()
	}

	return config.NewGetHostRuleOK().
		WithPayload(hostRules)
}

// RedirectHandler is handler for Redirect
func (c *Controller) RedirectHandler(params redirect.RedirectParams) middleware.Responder {
	hostRules, err := c.store.GetHostRules(params.Host)

	if nil != err {
		serverError := models.ServerError{Message: err.Error()}
		return redirect.NewRedirectInternalServerError().
			WithPayload(&serverError)
	}

	if nil == hostRules {
		return redirect.NewRedirectNotFound()
	}

	target := c.resolver.Resolve(*hostRules, params.SourcePath)

	return redirect.NewRedirectDefault(int(target.HTTPCode)).
		WithLocation(target.Path)
}
