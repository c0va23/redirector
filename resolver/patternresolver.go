package resolver

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/c0va23/redirector/log"
	"github.com/c0va23/redirector/models"
)

var placeholderRegexp = regexp.MustCompile("{(\\d+)}")

var patternLogger = log.NewLogger("PatterResolver", logrus.InfoLevel)

// PatternResolver resolve pathes with patterns and replace values in placeholders
func PatternResolver(
	rule models.Rule,
	sourcePath string,
) *models.Target {
	pattern, err := regexp.Compile(rule.SourcePath)
	if nil != err {
		patternLogger.Errorf("Error with path patter %s: %v", rule.SourcePath, err)
		return nil
	}

	matches := pattern.FindAllStringSubmatch(sourcePath, -1)
	patternLogger.Infof("Matches: %+v", matches)
	switch len(matches) {
	case 0:
		return nil
	case 1:
		target := buildTargetWithPlaceholders(matches[0][1:], rule.Target)
		return &target
	default:
		patternLogger.Infof(`Pattern "%s" match more one times path "%s"`, rule.SourcePath, sourcePath)
		return nil
	}
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
