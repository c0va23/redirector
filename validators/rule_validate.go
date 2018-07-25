package validators

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/c0va23/redirector/models"
	"github.com/c0va23/redirector/resolvers"
)

// ValidateRule return nil if rule is valid or models.ModelValidationError if
// is invalid
func ValidateRule(rule models.Rule) (
	modelError models.ModelValidationError,
	valid bool,
) {
	if models.RuleResolverSimple != rule.Resolver &&
		models.RuleResolverPattern != rule.Resolver {
		modelError = addFieldError(modelError, "resolver", "rule.resolver.unknown")
	}

	if "" == rule.SourcePath {
		modelError = addFieldError(modelError, "sourcePath", "required")
	}

	if models.RuleResolverPattern == rule.Resolver {
		modelError = validatePattern(modelError, rule)
	}

	if targetError, valid := ValidateTarget(rule.Target); !valid {
		modelError = addEmbedError(modelError, "target", targetError)
	}

	valid = isEmptyModelError(modelError)
	return
}

func validatePattern(
	modelError models.ModelValidationError,
	rule models.Rule,
) models.ModelValidationError {
	pattern, errPattern := regexp.Compile(rule.SourcePath)
	if nil != errPattern {
		return addFieldError(modelError, "sourcePath", "rule.sourcePath.invalidPattern")
	}

	placeholders := resolvers.PlaceholderRegexp.FindAllString(rule.Target.Path, -1)
	if len(placeholders) != pattern.NumSubexp() {
		return addFieldError(modelError, "target.path", "target.path.missedPlaceholder")
	}

	sort.StringSlice(placeholders).Sort()
	for index, placeholder := range placeholders {
		if placeholder != fmt.Sprintf("{%d}", index) {
			return addFieldError(modelError, "target.path", "target.path.invalidPlaceholderIndex")
		}
	}

	return modelError
}
