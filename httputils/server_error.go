package httputils

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/sirupsen/logrus"
)

// BuildServerErrorHandler build error handler
func BuildServerErrorHandler(
	redirectHandler http.Handler,
	logger *logrus.Logger,
) func(http.ResponseWriter, *http.Request, error) {
	return func(
		rw http.ResponseWriter,
		req *http.Request,
		err error,
	) {
		if apiErr, ok := err.(errors.Error); ok && http.StatusNotFound == apiErr.Code() {
			redirectHandler.ServeHTTP(rw, req)
		} else {
			logger.Errorf("ServerError %#v", err)
			errors.ServeError(rw, req, err)
		}
	}
}
