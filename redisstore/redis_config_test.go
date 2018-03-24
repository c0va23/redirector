package redisstore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/c0va23/redirector/redisstore"
)

func TestParseRedisURI_SuccessFull(t *testing.T) {
	a := assert.New(t)
	redisConfig, err := redisstore.ParseRedisURI("redis://localhost:16379/3")
	a.Nil(err)

	a.Equal(
		&redisstore.RedisConfig{
			Host:     "localhost",
			Port:     16379,
			Database: 3,
		},
		redisConfig,
	)
}
func TestParseRedisURI_SuccessDefaults(t *testing.T) {
	a := assert.New(t)
	redisConfig, err := redisstore.ParseRedisURI("redis://localhost")
	a.Nil(err)

	a.Equal(
		&redisstore.RedisConfig{
			Host:     "localhost",
			Port:     redisstore.DefaultPort,
			Database: redisstore.DefaultDatabase,
		},
		redisConfig,
	)
}

func TestParseRedisURI_InvalidURI(t *testing.T) {
	a := assert.New(t)
	_, err := redisstore.ParseRedisURI("")
	a.EqualError(err, "URI error: parse : empty url")
}

func TestParseRedisURI_InvalidScheme(t *testing.T) {
	a := assert.New(t)
	_, err := redisstore.ParseRedisURI("tcp://localhost")
	a.EqualError(err, "Unexpected scheme: tcp")
}

func TestParseRedisURI_InvalidPort(t *testing.T) {
	a := assert.New(t)
	_, err := redisstore.ParseRedisURI("redis://localhost:123456")
	a.EqualError(err, "Port error: strconv.ParseUint: parsing \"123456\": value out of range")
}

func TestParseRedisURI_InvalidDatabase(t *testing.T) {
	a := assert.New(t)
	_, err := redisstore.ParseRedisURI("redis://localhost/")
	a.EqualError(err, "Parse DB number error: strconv.ParseUint: parsing \"\": invalid syntax")
}
