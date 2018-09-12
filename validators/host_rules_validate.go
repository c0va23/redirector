package validators

import (
	"fmt"
	"regexp"

	"github.com/c0va23/redirector/models"
)

var hostRegex = regexp.MustCompile("^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$")

const maxHostLength = 255

const (
	hostRulesHostField = "host"
)

// ValidateHostRules with all embedded structs
func ValidateHostRules(hostRules models.HostRules) (
	modelError models.ModelValidationError,
	valid bool,
) {
	modelError = validateHost(modelError, hostRules.Host)

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

func validateHost(modelError models.ModelValidationError, host string) models.ModelValidationError {
	if "" == host {
		modelError = addFieldError(modelError, hostRulesHostField, "required")
	}

	if len(host) > maxHostLength {
		modelError = addFieldError(modelError, hostRulesHostField, "hostRules.host.tooLong")
	}

	if !hostRegex.MatchString(host) {
		modelError = addFieldError(modelError, hostRulesHostField, "hostRules.host.invalidPattern")
	}

	return modelError
}
