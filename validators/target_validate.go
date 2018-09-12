package validators

import (
	"github.com/c0va23/redirector/models"
)

const minHTTPCode = 300
const maxHTTPCode = 399

// ValidateTarget return validation error list and valid true if invalid.
// Otherwise return nil list and valid false.
func ValidateTarget(target models.Target) (
	modelError models.ModelValidationError,
	valid bool,
) {
	if !(minHTTPCode <= target.HTTPCode && target.HTTPCode <= maxHTTPCode) {
		modelError = addFieldError(modelError, "httpCode", "target.httpCode.outOfRange")
	}

	if "" == target.Path {
		modelError = addFieldError(modelError, "path", "required")
	}

	valid = isEmptyModelError(modelError)

	return
}
