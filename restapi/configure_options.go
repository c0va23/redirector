package restapi

import (
	"log"

	"github.com/go-openapi/swag"

	"github.com/c0va23/redirector/memstore"
	"github.com/c0va23/redirector/resolver"
	"github.com/c0va23/redirector/restapi/operations"
	"github.com/c0va23/redirector/store"
)

var appOptions struct {
	StoreType string `short:"s" long:"store-type" description:"Type of store engine" choice:"memory" choice:"redis" default:"memory"`
	RedisURI  string `long:"redis-uri" description:"Connection URI for Redis. Required if store-type equal redis"`
}

func configureFlags(api *operations.RedirectorAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		swag.CommandLineOptionsGroup{
			ShortDescription: "store",
			Options:          &appOptions,
		},
	}
}

func buildStore() store.Store {
	switch appOptions.StoreType {
	case "memory":
		return memstore.NewMemStore()
	default:
		log.Panicf("Unknown store type: %s", appOptions.StoreType)
		return nil
	}
}

func buildResolver() resolver.Resolver {
	return resolver.NewSimpleResolver()
}
