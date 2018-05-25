package resolver

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/c0va23/redirector/models"
)

var placeholderRegexp = regexp.MustCompile("{(\\d+)}")

// PatternResolver resolve pathes with patterns and replace values in placeholders
type PatternResolver struct{}

// Resolve implement Resolver.Resolve
func (r *PatternResolver) Resolve(
	hostRules models.HostRules,
	sourcePath string,
) models.Target {
	log.Printf("sourcePath: %s", sourcePath)
	for _, rule := range hostRules.Rules {
		log.Printf("Pattern: %s", rule.SourcePath)
		pattern, err := regexp.Compile(rule.SourcePath)
		if nil != err {
			log.Printf("Error with path patter %s: %v", rule.SourcePath, err)
			continue
		}

		matches := pattern.FindAllStringSubmatch(sourcePath, -1)
		log.Printf("Matches: %+v", matches)
		switch len(matches) {
		case 0:
			continue
		case 1:
			return buildTargetWithPlaceholders(matches[0][1:], rule.Target)
		default:
			log.Printf(`Pattern "%s" match more one times path "%s"`, rule.SourcePath, sourcePath)
		}
	}

	return hostRules.DefaultTarget
}

func buildTargetWithPlaceholders(matches []string, target models.Target) models.Target {
	targetPath := target.Path
	for index, value := range matches {
		placeholder := fmt.Sprintf("{%d}", index)
		targetPath = strings.Replace(targetPath, placeholder, value, -1)
	}

	return models.Target{
		HTTPCode: target.HTTPCode,
		Path:     targetPath,
	}
}
