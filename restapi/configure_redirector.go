// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"
	"strings"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	cors "github.com/rs/cors"
	graceful "github.com/tylerb/graceful"

	"github.com/c0va23/redirector/controllers"
	"github.com/c0va23/redirector/restapi/operations"
	"github.com/c0va23/redirector/restapi/operations/config"
	"github.com/c0va23/redirector/restapi/operations/redirect"
)

//go:generate swagger generate server --target .. --name  --spec ../api.yml

func configureAPI(api *operations.RedirectorAPI) http.Handler {
	store := buildStore()
	resolver := buildResolver()

	controller := controllers.NewController(store, resolver)

	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.APISecurityAuth = basicAuth

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.ConfigListHostRulesHandler = config.ListHostRulesHandlerFunc(controller.ListHostRulesHandler)
	api.ConfigCreateHostRulesHandler = config.CreateHostRulesHandlerFunc(controller.CreateHostRulesHandler)
	api.ConfigUpdateHostRulesHandler = config.UpdateHostRulesHandlerFunc(controller.UpdateHostRulesHandler)
	api.ConfigGetHostRuleHandler = config.GetHostRuleHandlerFunc(controller.GetHostRulesHandler)
	api.ConfigDeleteHostRulesHandler = config.DeleteHostRulesHandlerFunc(controller.DeleteHostRulesHandler)
	api.RedirectRedirectHandler = redirect.RedirectHandlerFunc(controller.RedirectHandler)

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme, addr string) {
}

func returnHostHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		hostParts := strings.Split(req.Host, ":")
		host := hostParts[0]
		req.Header.Add("Host", host)
		next.ServeHTTP(res, req)
	})
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return returnHostHeader(handler)
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return cors.AllowAll().Handler(handler)
}
