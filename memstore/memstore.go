package memstore

import (
	"sync"

	"github.com/c0va23/redirector/models"
)

// MemStore is in-memory implementation of store.Store
type MemStore struct {
	sync.RWMutex
	listHostRules []models.HostRules
}

// NewMemStore create new MemStore
func NewMemStore() *MemStore {
	return &MemStore{
		RWMutex:       sync.RWMutex{},
		listHostRules: []models.HostRules{},
	}
}

// ListHostRules return list of host rules from MemStore
func (memStore *MemStore) ListHostRules() ([]models.HostRules, error) {
	memStore.RLock()
	defer memStore.RUnlock()
	return memStore.listHostRules, nil
}

// ReplaceHostRules replace or create host rule into MemStore
func (memStore *MemStore) ReplaceHostRules(newHostRules models.HostRules) error {
	memStore.Lock()
	defer memStore.Unlock()

	updated := false
	for index, hostRules := range memStore.listHostRules {
		if newHostRules.Host == hostRules.Host {
			memStore.listHostRules[index] = newHostRules
			updated = true
		}
	}

	if !updated {
		memStore.listHostRules = append(memStore.listHostRules, newHostRules)
	}
	return nil
}

// GetHostRules return HostRule by host
func (memStore *MemStore) GetHostRules(host string) (*models.HostRules, error) {
	for _, hostRules := range memStore.listHostRules {
		if host == hostRules.Host {
			return &hostRules, nil
		}
	}

	return nil, nil
}
