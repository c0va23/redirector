package handlers

import (
	"sort"

	"github.com/go-openapi/runtime/middleware"

	"github.com/c0va23/redirector/models"
	"github.com/c0va23/redirector/restapi/operations/config"
	"github.com/c0va23/redirector/restapi/operations/redirect"
	"github.com/c0va23/redirector/store"
)

// ConfigHandlers implement methods into restapi
type ConfigHandlers struct {
	store store.Store
}

// NewConfigHandlers initialize new controller
func NewConfigHandlers(store store.Store) ConfigHandlers {
	return ConfigHandlers{
		store: store,
	}
}

// ListHostRulesHandler is handler for ListHostRules
func (configHandler *ConfigHandlers) ListHostRulesHandler(
	params config.ListHostRulesParams,
	_principal interface{},
) middleware.Responder {
	listHostRules, err := configHandler.store.ListHostRules()

	switch err {
	case nil:
		sort.Slice(listHostRules, func(i, j int) bool {
			return listHostRules[i].Host < listHostRules[j].Host
		})
		return config.NewListHostRulesOK().
			WithPayload(listHostRules)
	default:
		return config.
			NewListHostRulesInternalServerError().
			WithPayload(&models.ServerError{Message: err.Error()})
	}
}

// CreateHostRulesHandler is handler for CreateHostRules
func (configHandler *ConfigHandlers) CreateHostRulesHandler(
	params config.CreateHostRulesParams,
	_principal interface{},
) middleware.Responder {
	err := configHandler.store.CreateHostRules(params.HostRules)

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
func (configHandler *ConfigHandlers) UpdateHostRulesHandler(
	params config.UpdateHostRulesParams,
	_principal interface{},
) middleware.Responder {
	err := configHandler.store.UpdateHostRules(params.Host, params.HostRules)

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
func (configHandler *ConfigHandlers) GetHostRulesHandler(
	params config.GetHostRuleParams,
	_principal interface{},
) middleware.Responder {
	hostRules, err := configHandler.store.GetHostRules(params.Host)

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
func (configHandler *ConfigHandlers) DeleteHostRulesHandler(
	params config.DeleteHostRulesParams,
	_principal interface{},
) middleware.Responder {
	err := configHandler.store.DeleteHostRules(params.Host)
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
func (configHandler *ConfigHandlers) HealthCheckHandler(
	params redirect.HealthcheckParams,
) middleware.Responder {
	err := configHandler.store.CheckHealth()

	switch err {
	case nil:
		return redirect.NewHealthcheckOK()
	default:
		return redirect.NewHealthcheckInternalServerError()
	}
}
