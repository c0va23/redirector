package controllers

import (
	"net/http"
	"strings"

	"github.com/c0va23/redirector/log"
	"github.com/c0va23/redirector/resolver"
	"github.com/c0va23/redirector/store"
	"github.com/sirupsen/logrus"
)

const locationHeader = "Location"

const (
	notFoundMessage            = "Not found"
	internalServerErrorMessage = "Internal server error"
)

var redirectLogger = log.NewLogger("RedirectHandler", logrus.InfoLevel)

// RedirectHandler is handler for redirect requests
type RedirectHandler struct {
	store.Store
	resolver.HostRulesResolver
}

// NewRedirectHandler build new RedirectHandler
func NewRedirectHandler(
	store store.Store,
	resolver resolver.HostRulesResolver,
) *RedirectHandler {
	return &RedirectHandler{
		Store:             store,
		HostRulesResolver: resolver,
	}
}

// ServeHTTP handle all requests for redirect.
// If target not found, then respond with 404 code.
// If store return error, then respond with 500 code.
func (rh *RedirectHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	hostParts := strings.Split(req.Host, ":")
	host := hostParts[0]

	hostRules, err := rh.GetHostRules(host)

	switch err {
	case nil:
		target := rh.Resolve(*hostRules, req.RequestURI)
		rw.Header().Add(locationHeader, target.Path)
		rw.WriteHeader(int(target.HTTPCode))
	case store.ErrNotFound:
		rw.WriteHeader(http.StatusNotFound)
		writeBody(rw, notFoundMessage)
	default:
		redirectLogger.WithError(err).Errorf("Redirect error: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		writeBody(rw, internalServerErrorMessage)
	}
}

func writeBody(rw http.ResponseWriter, message string) {
	_, err := rw.Write(([]byte)(message))
	if nil != err {
		redirectLogger.
			WithError(err).
			Errorf("Error on write not found body: %s", err)
	}
}
