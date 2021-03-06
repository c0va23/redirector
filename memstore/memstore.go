package memstore

import (
	"sync"

	"github.com/c0va23/redirector/models"
	"github.com/c0va23/redirector/store"
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

// CreateHostRules create new HostRules if it not exists
func (memStore *MemStore) CreateHostRules(newHostRules models.HostRules) error {
	memStore.Lock()
	defer memStore.Unlock()

	for _, hostRules := range memStore.listHostRules {
		if newHostRules.Host == hostRules.Host {
			return store.ErrExists
		}
	}

	memStore.listHostRules = append(memStore.listHostRules, newHostRules)

	return nil
}

// UpdateHostRules is update host rules if exists
func (memStore *MemStore) UpdateHostRules(host string, updatedHostRules models.HostRules) error {
	memStore.Lock()
	defer memStore.Unlock()

	for index, hostRules := range memStore.listHostRules {
		if updatedHostRules.Host == hostRules.Host {
			memStore.listHostRules[index] = updatedHostRules
			return nil
		}
	}

	return store.ErrNotFound
}

// GetHostRules return HostRule by host
func (memStore *MemStore) GetHostRules(host string) (*models.HostRules, error) {
	for _, hostRules := range memStore.listHostRules {
		if host == hostRules.Host {
			return &hostRules, nil
		}
	}

	return nil, store.ErrNotFound
}

// DeleteHostRules delete host rules by host
func (memStore *MemStore) DeleteHostRules(host string) error {
	memStore.Lock()
	defer memStore.Unlock()

	for index, hostRules := range memStore.listHostRules {
		if host == hostRules.Host {
			memStore.listHostRules = append(
				memStore.listHostRules[:index],
				memStore.listHostRules[index+1:]...,
			)
			return nil
		}
	}
	return store.ErrNotFound
}

// CheckHealth always return nil
func (memStore *MemStore) CheckHealth() error {
	return nil
}
