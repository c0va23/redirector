package handlers

import (
	"github.com/c0va23/redirector/models"
	"github.com/c0va23/redirector/restapi/operations/config"

	"github.com/go-openapi/runtime/middleware"
)

// LocalesHandler is handler for locales
type LocalesHandler models.Locales

// Handle implement config.LocalesHandler.handle
func (locales LocalesHandler) Handle(config.LocalesParams) middleware.Responder {
	return config.
		NewLocalesOK().
		WithPayload(models.Locales(locales))
}
