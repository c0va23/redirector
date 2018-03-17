package controllers

import (
	"log"

	"github.com/go-openapi/runtime/middleware"

	"github.com/c0va23/redirector/models"
	"github.com/c0va23/redirector/restapi/operations"
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
func (c *Controller) ListHostRulesHandler(params operations.ListHostRulesParams) middleware.Responder {
	hostRules, err := c.store.ListHostRules()
	if nil != err {
		serverError := models.ServerError{Message: err.Error()}
		return operations.NewListHostRulesInternalServerError().WithPayload(&serverError)
	}

	log.Printf("%+v", hostRules)
	listHostRules := make([]*models.HostRule, 0, len(hostRules))

	for _, hostRule := range hostRules {
		listHostRules = append(listHostRules, &hostRule)
	}
	return operations.NewListHostRulesOK().WithPayload(listHostRules)
}

// ReplaceHostRulesHandler is handler for ReplaceHostRules
func (c *Controller) ReplaceHostRulesHandler(params operations.ReplaceHostRuleParams) middleware.Responder {
	err := c.store.ReplaceHostRule(*params.HostRule)

	if nil != err {
		serverError := models.ServerError{Message: err.Error()}
		return operations.NewReplaceHostRuleInternalServerError().WithPayload(&serverError)
	}

	return operations.NewReplaceHostRuleOK().
		WithPayload(params.HostRule)
}
