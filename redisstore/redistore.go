package redisstore

import (
	"encoding/json"
	"strings"

	"github.com/mediocregopher/radix.v2/util"

	"github.com/c0va23/redirector/models"
)

// RedisStore implement store with Redis backend
type RedisStore struct {
	util.Cmder
}

// NewRedisStore create new RedisStore
func NewRedisStore(cmder util.Cmder) RedisStore {
	return RedisStore{
		Cmder: cmder,
	}
}

// ListHostRules implement Store.ListHostRules
func (rs *RedisStore) ListHostRules() ([]models.HostRules, error) {
	scanner := util.NewScanner(rs.Cmder, util.ScanOpts{Command: "SCAN"})
	hosts := []string{}
	for scanner.HasNext() {
		host := scanner.Next()
		hosts = append(hosts, host)
	}
	if err := scanner.Err(); nil != err {
		return nil, err
	}

	listHostRules := make([]models.HostRules, 0, len(hosts))
	for _, host := range hosts {
		hostRuleJSON, err := rs.Cmder.Cmd("GET", host).Str()
		if nil != err {
			return nil, err
		}

		decoder := json.NewDecoder(strings.NewReader(hostRuleJSON))
		var hostRules models.HostRules
		if err := decoder.Decode(&hostRules); nil != err {
			return nil, err
		}
		listHostRules = append(listHostRules, hostRules)
	}

	return listHostRules, nil
}

// ReplaceHostRules implement Store.ReplaceHostRules
func (rs *RedisStore) ReplaceHostRules(hostRules models.HostRules) error {
	json, _ := json.Marshal(hostRules)
	resp := rs.Cmder.Cmd("SET", hostRules.Host, string(json))

	if _, err := resp.Str(); nil != err {
		return err
	}

	return nil
}

// GetHostRules implement Store.GetHostRules
func (rs *RedisStore) GetHostRules(host string) (*models.HostRules, error) {
	return nil, nil
}
