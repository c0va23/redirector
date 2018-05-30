package resolvers

import (
	"github.com/c0va23/redirector/models"
)

// SimpleResolver is resolver source path by full match (without pattern)
func SimpleResolver(
	rule models.Rule,
	sourcePath string,
) *models.Target {
	if rule.SourcePath == sourcePath {
		return &rule.Target
	}

	return nil
}
