package restapi

import (
	"errors"

	"github.com/go-openapi/swag"
	"github.com/sirupsen/logrus"

	"github.com/c0va23/redirector/log"
	"github.com/c0va23/redirector/memstore"
	"github.com/c0va23/redirector/redisstore"
	"github.com/c0va23/redirector/resolver"
	"github.com/c0va23/redirector/restapi/operations"
	"github.com/c0va23/redirector/store"
)

var appOptions struct {
	StoreType     string `short:"s" long:"store-type" description:"Type of store engine" choice:"memory" choice:"redis" default:"memory"`
	RedisURI      string `long:"redis-uri" description:"Connection URI for Redis. Required if store-type equal redis"`
	RedisPoolSize int    `long:"redis-pool-size" description:"Redis pool size" default:"10"`
	BasicUsername string `long:"basic-username" short:"u" env:"BASIC_USERNAME" description:"Username for Basic auth" required:"true"`
	BasicPassword string `long:"basic-password" short:"p" env:"BASIC_PASSWORD" description:"Password for Basic auth." required:"true"`
	Resolver      string `long:"resolver" env:"RESOLVER" description:"Simple or Pattern" default:"simple"`
}

var builderLogger = log.NewLogger("builer", logrus.InfoLevel)

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
	case "redis":
		client, err := redisstore.BuildRedisPool(appOptions.RedisURI, appOptions.RedisPoolSize)
		if nil != err {
			builderLogger.Fatalf("Redis error: %s", err)
		}
		return redisstore.NewRedisStore(client)
	default:
		builderLogger.Fatalf("Unknown store type: %s", appOptions.StoreType)
		return nil
	}
}

func buildResolver() resolver.Resolver {
	switch appOptions.Resolver {
	case "simple":
		return new(resolver.SimpleResolver)
	case "pattern":
		return new(resolver.PatternResolver)
	default:
		builderLogger.Fatalf("Unknown resolver: %s", appOptions.Resolver)
		return nil
	}
}

var errInvalidBasicCredentials = errors.New("Invalid Basic credentials")

var authLogger = log.NewLogger("auth", logrus.InfoLevel)

func basicAuth(userename, password string) (interface{}, error) {
	if appOptions.BasicUsername != userename || appOptions.BasicPassword != password {
		authLogger.WithError(errInvalidBasicCredentials).
			Errorf(`Invalid credential for username "%s"`, userename)
		return false, errInvalidBasicCredentials
	}

	return true, nil
}
