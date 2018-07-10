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

// ErrInvalidBasicCredentials is error when user pass invalid username or password
var ErrInvalidBasicCredentials = errors.New(
	http.StatusUnauthorized,
	"Invalid Basic credentials",
)

// BuildBasicAuth build basicAuth checker
func BuildBasicAuth(
	expectedUsername string,
	expectedPassword string,
	logger *logrus.Logger,
) func(username, password string) (interface{}, error) {
	return func(userename, password string) (interface{}, error) {
		if expectedUsername != userename || expectedPassword != password {
			logger.WithError(ErrInvalidBasicCredentials).
				Warnf(`Invalid credential for username "%s"`, userename)
			return false, ErrInvalidBasicCredentials
		}

		return true, nil
	}
}
