package store

import (
	"github.com/c0va23/redirector/models"
)

// Store is interface for stores
type Store interface {
	// ListHostRules return list of models.HostRules or error
	ListHostRules() ([]models.HostRule, error)
	// ReplaceHostRule replace HostRule if it exists or add new HostRule if it not exists
	ReplaceHostRule(models.HostRule) error
	// GetHostRules return HostRules by Host if it exists.
	// Or return nil if Host not exists. Otherwise return error.
	GetHostRules(string) (*models.HostRule, error)
}
