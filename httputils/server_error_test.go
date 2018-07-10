package httputils_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/c0va23/redirector/httputils"

	"github.com/go-openapi/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type HTTPHandlerMock struct {
	mock.Mock
}

func (h *HTTPHandlerMock) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.Mock.MethodCalled("ServeHTTP", rw, r)
}

func TestBuildServerErrorHandler_Redirect(t *testing.T) {
	rw := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "https://test.org", nil)

	handlerMock := new(HTTPHandlerMock)
	handlerMock.On("ServeHTTP", rw, r)

	logger := logrus.New()
	handler := httputils.BuildServerErrorHandler(handlerMock, logger)

	err := errors.New(http.StatusNotFound, "NotFound")
	handler(rw, r, err)

	handlerMock.AssertExpectations(t)
}

func TestBuildServerErrorHandler_Error(t *testing.T) {
	rw := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "https://test.org", nil)

	handlerMock := new(HTTPHandlerMock)
	logger := logrus.New()
	handler := httputils.BuildServerErrorHandler(handlerMock, logger)

	err := errors.New(http.StatusUnauthorized, "Unauthorized")
	handler(rw, r, err)

	a := assert.New(t)

	a.Equal(int(err.Code()), rw.Code)
}
