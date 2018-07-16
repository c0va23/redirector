package handlers_test

import (
	"testing"

	"github.com/c0va23/redirector/handlers"
	"github.com/c0va23/redirector/models"
	"github.com/c0va23/redirector/restapi/operations/config"

	"github.com/stretchr/testify/assert"
)

func TestLocalesHandler(t *testing.T) {
	a := assert.New(t)

	locales := models.Locales{}
	localesHandler := handlers.LocalesHandler(locales)

	a.Implements((*config.LocalesHandler)(nil), &localesHandler)
}

func TestLocalesHandler_Handle(t *testing.T) {
	a := assert.New(t)

	locales := models.Locales{}
	localesHandler := handlers.LocalesHandler(locales)

	a.Equal(
		config.NewLocalesOK().
			WithPayload(locales),
		localesHandler.Handle(config.LocalesParams{}),
	)
}
