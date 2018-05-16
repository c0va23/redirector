package redisstore

import (
	"encoding/json"

	"github.com/mediocregopher/radix.v2/redis"
	"github.com/mediocregopher/radix.v2/util"

	"github.com/c0va23/redirector/models"
	"github.com/c0va23/redirector/store"
)

// RedisStore implement store with Redis backend
type RedisStore struct {
	util.Cmder
}

// NewRedisStore create new RedisStore
func NewRedisStore(cmder util.Cmder) *RedisStore {
	return &RedisStore{
		Cmder: cmder,
	}
}

func (rs *RedisStore) listHosts() ([]string, error) {
	scanner := util.NewScanner(rs.Cmder, util.ScanOpts{Command: "SCAN"})

	hosts := []string{}
	for scanner.HasNext() {
		host := scanner.Next()
		hosts = append(hosts, host)
	}

	if err := scanner.Err(); nil != err {
		return nil, err
	}

	return hosts, nil
}

func (rs *RedisStore) fetchHostRules(hosts []string) ([]models.HostRules, error) {
	listHostRules := make([]models.HostRules, 0, len(hosts))
	for _, host := range hosts {
		hostRuleJSON, err := rs.Cmder.Cmd("GET", host).Bytes()
		if nil != err {
			return nil, err
		}

		var hostRules models.HostRules
		if err := json.Unmarshal(hostRuleJSON, &hostRules); nil != err {
			return nil, err
		}

		listHostRules = append(listHostRules, hostRules)
	}

	return listHostRules, nil
}

// ListHostRules implement Store.ListHostRules
func (rs *RedisStore) ListHostRules() ([]models.HostRules, error) {
	hosts, err := rs.listHosts()
	if nil != err {
		return nil, err
	}

	listHostRules, err := rs.fetchHostRules(hosts)
	if nil != err {
		return nil, err
	}

	return listHostRules, nil
}

// GetHostRules implement Store.GetHostRules
func (rs *RedisStore) GetHostRules(host string) (*models.HostRules, error) {
	resp := rs.Cmder.Cmd("GET", host)

	value, err := resp.Bytes()
	if redis.ErrRespNil == err {
		return nil, nil
	} else if nil != err {
		return nil, err
	}

	var hostRules models.HostRules
	if err := json.Unmarshal(value, &hostRules); nil != err {
		return nil, err
	}

	return &hostRules, nil
}

// CreateHostRules create host rules if not exists
func (rs *RedisStore) CreateHostRules(hostRules models.HostRules) error {
	json, _ := json.Marshal(hostRules)
	resp := rs.Cmder.Cmd("SETNX", hostRules.Host, string(json))

	if code, err := resp.Int(); nil != err {
		return err
	} else if 0 == code {
		return store.ErrExists
	}

	return nil
}

// UpdateHostRules update host rules if it exists
func (rs *RedisStore) UpdateHostRules(host string, hostRules models.HostRules) error {
	if code, err := rs.Cmd("EXISTS", host).Int(); nil != err {
		return err
	} else if 0 == code {
		return store.ErrNotFound
	}

	if host != hostRules.Host {
		if code, err := rs.Cmd("EXISTS", hostRules.Host).Int(); nil != err {
			return err
		} else if 1 == code {
			return store.ErrExists
		}
	}

	hostRulesJSON, _ := json.Marshal(&hostRules)
	setResp := rs.Cmder.Cmd("SET", hostRules.Host, string(hostRulesJSON))

	if _, err := setResp.Str(); nil != err {
		return err
	}

	if host != hostRules.Host {
		if _, err := rs.Cmd("DEL", host).Int(); nil != err {
			return err
		}
	}

	return nil
}

// DeleteHostRules if host exists
func (rs *RedisStore) DeleteHostRules(host string) error {
	if count, err := rs.Cmder.Cmd("DEL", host).Int(); nil != err {
		return err
	} else if 0 == count {
		return store.ErrNotFound
	}

	return nil
}
