// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	runtime "github.com/go-openapi/runtime"
	cors "github.com/rs/cors"

	"github.com/c0va23/redirector/handlers"
	"github.com/c0va23/redirector/httputils"
	"github.com/c0va23/redirector/log"
	"github.com/c0va23/redirector/resolvers"
	"github.com/c0va23/redirector/restapi/operations"
	"github.com/c0va23/redirector/restapi/operations/config"
	"github.com/c0va23/redirector/restapi/operations/redirect"
)

//go:generate swagger generate server --target .. --name  --spec ../api.yml

var configLogger = log.NewLeveledLogger("config")

func configureAPI(api *operations.RedirectorAPI) http.Handler {
	store := buildStore()

	configHandlers := handlers.NewConfigHandlers(store)

	resolver := resolvers.MultiHostRulesResolver(resolvers.DefaultResolvers)
	redirectHandler := handlers.NewRedirectHandler(store, resolver)

	// configure the api here
	api.ServeError = httputils.BuildServerErrorHandler(
		redirectHandler,
		configLogger,
	)

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	api.Logger = configLogger.Infof

	api.APISecurityAuth = basicAuth

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.ConfigListHostRulesHandler = config.ListHostRulesHandlerFunc(configHandlers.ListHostRulesHandler)
	api.ConfigCreateHostRulesHandler = config.CreateHostRulesHandlerFunc(configHandlers.CreateHostRulesHandler)
	api.ConfigUpdateHostRulesHandler = config.UpdateHostRulesHandlerFunc(configHandlers.UpdateHostRulesHandler)
	api.ConfigGetHostRuleHandler = config.GetHostRuleHandlerFunc(configHandlers.GetHostRulesHandler)
	api.ConfigDeleteHostRulesHandler = config.DeleteHostRulesHandlerFunc(configHandlers.DeleteHostRulesHandler)
	api.RedirectHealthcheckHandler = redirect.HealthcheckHandlerFunc(configHandlers.HealthCheckHandler)

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
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	corsHandler := cors.AllowAll().Handler(handler)
	return log.Request(corsHandler)
}
