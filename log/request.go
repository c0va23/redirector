package log

import (
	"net/http"
	"time"

	logrus "github.com/sirupsen/logrus"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	Status int
}

func newResponseWriterWrapper(
	responseWriter http.ResponseWriter,
) *responseWriterWrapper {
	return &responseWriterWrapper{
		ResponseWriter: responseWriter,
	}
}

func (rww *responseWriterWrapper) WriteHeader(status int) {
	rww.Status = status
	rww.ResponseWriter.WriteHeader(status)
}

func (rww *responseWriterWrapper) Write(data []byte) (int, error) {
	return rww.ResponseWriter.Write(data)
}

func (rww *responseWriterWrapper) Header() http.Header {
	return rww.ResponseWriter.Header()
}

var requestLogger = NewLogger("HTTP", logrus.InfoLevel)

// Request is middleware for request logging
func Request(next http.Handler) http.Handler {
	return http.HandlerFunc(func(
		res http.ResponseWriter,
		req *http.Request,
	) {
		resWrapper := newResponseWriterWrapper(res)
		startTime := time.Now()

		next.ServeHTTP(resWrapper, req)

		endTime := time.Now()

		entry := requestLogger.WithFields(map[string]interface{}{
			"method":   req.Method,
			"path":     req.RequestURI,
			"status":   resWrapper.Status,
			"duration": endTime.Sub(startTime).Seconds(),
		})

		if location := resWrapper.Header().Get("Location"); "" != location {
			entry = entry.WithField("location", location)
		}

		entry.Infof(
			"%s %s",
			req.Method,
			req.RequestURI,
		)
	})
}
