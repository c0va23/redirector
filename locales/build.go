package locales

//go:generate parcello

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/phogolabs/parcello"

	"github.com/c0va23/redirector/models"
)

const rootDir = "/"
const tomlExt = ".toml"

type group string
type key string
type localeFileStruct map[group]map[key]string

// DefaultLocale used for mark default locale
const DefaultLocale = "en"

func buildLocaleTranslations(
	translationMap localeFileStruct,
	localeCode string,
) models.LocaleTranslations {
	translations := []models.Translation{}

	for section, messagesMap := range translationMap {
		for messageKey, message := range messagesMap {
			translations = append(translations, models.Translation{
				Key:     fmt.Sprintf("%s.%s", section, messageKey),
				Message: message,
			})
		}
	}

	return models.LocaleTranslations{
		Code:          localeCode,
		Translations:  translations,
		DefaultLocale: DefaultLocale == localeCode,
	}
}

func parseTranslationFile(path string) (models.LocaleTranslations, error) {
	file, _ := parcello.Open(path)

	translationMap := localeFileStruct{}

	if _, err := toml.DecodeReader(file, &translationMap); nil != err {
		return models.LocaleTranslations{}, err
	}

	localeCode := strings.Replace(path[1:], tomlExt, "", 1)

	return buildLocaleTranslations(translationMap, localeCode), nil
}

func isTomlFile(info os.FileInfo) bool {
	return !info.IsDir() && tomlExt == filepath.Ext(info.Name())
}

// BuildLocales read locales and build Translations
func BuildLocales() (models.Locales, error) {
	localeList := models.Locales{}
	err := parcello.Manager.Walk(
		"/",
		func(path string, info os.FileInfo, _ error) error {
			if !isTomlFile(info) {
				return nil
			}

			localeTranslations, err := parseTranslationFile(path)
			if nil != err {
				return err
			}

			localeList = append(localeList, localeTranslations)
			return nil
		},
	)
	return localeList, err
}
