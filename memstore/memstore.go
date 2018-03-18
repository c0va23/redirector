package memstore

import (
	"github.com/c0va23/redirector/models"
	"log"
	"sync"
)

// MemStore is in-memory implementation of store.Store
type MemStore struct {
	sync.RWMutex
	hostRules []models.HostRule
}

// NewMemStore create new MemStore
func NewMemStore() MemStore {
	return MemStore{
		RWMutex:   sync.RWMutex{},
		hostRules: []models.HostRule{},
	}
}

// ListHostRules return list of host rules from MemStore
func (memStore *MemStore) ListHostRules() ([]models.HostRule, error) {
	memStore.RLock()
	defer memStore.RUnlock()
	return memStore.hostRules, nil
}

// ReplaceHostRule replace or create host rule into MemStore
func (memStore *MemStore) ReplaceHostRule(newHostRule models.HostRule) error {
	memStore.Lock()
	defer memStore.Unlock()

	updated := false
	for _, hostRule := range memStore.hostRules {
		if newHostRule.Host == hostRule.Host {
			log.Printf("Update %+v => %+v", newHostRule, hostRule)
			hostRule = newHostRule
			updated = true
		}
	}

	if !updated {
		memStore.hostRules = append(memStore.hostRules, newHostRule)
		log.Printf("Append %+v", memStore.hostRules)
	}
	return nil
}

// GetHostRules return HostRule by host
func (memStore *MemStore) GetHostRules(host string) (*models.HostRule, error) {
	for _, hostRule := range memStore.hostRules {
		if host == hostRule.Host {
			return &hostRule, nil
		}
	}

	return nil, nil
}
