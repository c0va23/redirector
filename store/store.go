package store

import (
	"errors"

	"github.com/c0va23/redirector/models"
)

var NotFound = errors.New("Host rules not found")
var Exists = errors.New("Host rules already exists")

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
}
