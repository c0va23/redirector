package redisstore

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/mediocregopher/radix.v2/util"
)

// Port is type for Redis port number
type Port uint16

// Database is type for redis database number
type Database int

// RedisConfig is type for redis config
type RedisConfig struct {
	Host     string
	Port     Port
	Database Database
	PoolSize uint
}

// DefaultPort is Default value for Port
const DefaultPort = Port(6379)

// DefaultDatabase is Default value for Database
const DefaultDatabase = Database(0)

const redisScheme = "redis"

func parseRedisPort(portSrc string) (Port, error) {
	if "" == portSrc {
		return DefaultPort, nil
	}

	p, err := strconv.ParseUint(portSrc, 10, 16)
	if nil != err {
		return DefaultPort, fmt.Errorf("Port error: %s", err)
	}

	return Port(p), nil
}

func parseRedisDB(path string) (Database, error) {
	if 0 == len(path) {
		return DefaultDatabase, nil
	}

	db, err := strconv.ParseUint(path[1:], 10, 61)
	if nil != err {
		return DefaultDatabase, fmt.Errorf("Parse DB number error: %s", err)
	}
	return Database(db), nil
}

// ParseRedisURI parse RedisURI
func ParseRedisURI(redisURI string) (*RedisConfig, error) {
	uri, err := url.ParseRequestURI(redisURI)
	if nil != err {
		return nil, fmt.Errorf("URI error: %s", err)
	}

	if redisScheme != uri.Scheme {
		return nil, fmt.Errorf("Unexpected scheme: %s", uri.Scheme)
	}

	port, err := parseRedisPort(uri.Port())
	if nil != err {
		return nil, err
	}

	db, err := parseRedisDB(uri.Path)
	if nil != err {
		return nil, err
	}

	return &RedisConfig{
		Port:     Port(port),
		Host:     uri.Hostname(),
		Database: db,
	}, nil
}

func buildCustomDialer(db Database) pool.DialFunc {
	return func(network, addr string) (*redis.Client, error) {
		client, err := redis.Dial(network, addr)
		if nil != err {
			return nil, fmt.Errorf("Dial error: %s", err)
		}

		if _, err := client.Cmd("SELECT", db).Str(); nil != err {
			return nil, fmt.Errorf("Select DB error: %s", err)
		}

		return client, nil
	}
}

// BuildRedisPool build new redis pool
func BuildRedisPool(redisURI string, poolSize int) (util.Cmder, error) {
	config, err := ParseRedisURI(redisURI)
	if nil != err {
		return nil, err
	}

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	return pool.NewCustom(
		"tcp",
		addr,
		poolSize,
		buildCustomDialer(config.Database),
	)
}
