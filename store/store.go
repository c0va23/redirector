package store

import (
	"github.com/c0va23/redirector/models"
)

type Store interface {
	ListHostRules() ([]models.HostRule, error)
	ReplaceHostRule(models.HostRule) error
}
