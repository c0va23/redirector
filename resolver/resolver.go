package resolver

import (
	"github.com/c0va23/redirector/models"
)

// Resolver used for resolve host and path to HTTP code and new URL
type Resolver interface {
	Resolve(models.HostRules, string) models.Target
}
