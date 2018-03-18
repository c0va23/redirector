package controllers

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/c0va23/redirector/models"
	"github.com/c0va23/redirector/resolver"
	"github.com/c0va23/redirector/restapi/operations"
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
func (c *Controller) ListHostRulesHandler(params operations.ListHostRulesParams) middleware.Responder {
	hostRules, err := c.store.ListHostRules()

	if nil != err {
		serverError := models.ServerError{Message: err.Error()}
		return operations.NewListHostRulesInternalServerError().WithPayload(&serverError)
	}

	return operations.NewListHostRulesOK().
		WithPayload(hostRules)
}

// ReplaceHostRulesHandler is handler for ReplaceHostRules
func (c *Controller) ReplaceHostRulesHandler(params operations.ReplaceHostRuleParams) middleware.Responder {
	err := c.store.ReplaceHostRule(params.HostRule)

	if nil != err {
		serverError := models.ServerError{Message: err.Error()}
		return operations.NewReplaceHostRuleInternalServerError().WithPayload(&serverError)
	}

	return operations.NewReplaceHostRuleOK().
		WithPayload(&params.HostRule)
}

// RedirectHandler is handler for Redirect
func (c *Controller) RedirectHandler(params operations.RedirectParams) middleware.Responder {
	hostRules, err := c.store.GetHostRules(params.Host)

	if nil != err {
		serverError := models.ServerError{Message: err.Error()}
		return operations.NewRedirectInternalServerError().
			WithPayload(&serverError)
	}

	if nil == hostRules {
		return operations.NewRedirectNotFound()
	}

	target := c.resolver.Resolve(*hostRules, params.SourcePath)

	return operations.NewRedirectDefault(int(target.HTTPCode)).
		WithLocation(target.TargetPath)
}
