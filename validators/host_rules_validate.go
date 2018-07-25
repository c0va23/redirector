package validators

import (
	"fmt"

	"github.com/c0va23/redirector/models"
)

// ValidateHostRules with all embedded structs
func ValidateHostRules(hostRules models.HostRules) (
	modelError models.ModelValidationError,
	valid bool,
) {
	if "" == hostRules.Host {
		modelError = addFieldError(modelError, "host", "required")
	}

	if defaultTargetError, valid := ValidateTarget(hostRules.DefaultTarget); !valid {
		modelError = addEmbedError(modelError, "defaultTarget", defaultTargetError)
	}

	for index, rule := range hostRules.Rules {
		if ruleError, valid := ValidateRule(rule); !valid {
			modelError = addEmbedError(
				modelError,
				fmt.Sprintf("rules.%d", index),
				ruleError,
			)
		}
	}

	valid = isEmptyModelError(modelError)
	return
}
