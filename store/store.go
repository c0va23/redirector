package store

import (
	"errors"

	"github.com/c0va23/redirector/models"
)

// ErrNotFound is error for HostRules not found
var ErrNotFound = errors.New("Host rules not found")

// ErrExists is error when HostRules already exists
var ErrExists = errors.New("Host rules already exists")

// Store is interface for stores
type Store interface {
	// ListHostRules return list of models.HostRules or error
	ListHostRules() ([]models.HostRules, error)
	// GetHostRules return HostRules by Host if it exists.
	// Or return nil if Host not exists. Otherwise return error.
	GetHostRules(string) (*models.HostRules, error)
	// CreateHostRule create HostRule if it not exists
	CreateHostRules(models.HostRules) error
	// UpdateHostRules update HostRules if it exists
	UpdateHostRules(host string, hostRules models.HostRules) error
	// DeleteHostRules delete HostRules if it exists
	DeleteHostRules(host string) error
}
