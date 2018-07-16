package locales_test

import (
	"reflect"
	"testing"

	"github.com/c0va23/redirector/locales"
	"github.com/c0va23/redirector/models"

	"github.com/stretchr/testify/assert"
)

func TestBuildLocales_NotEmpty(t *testing.T) {
	a := assert.New(t)

	translations, err := locales.BuildLocales()

	a.Nil(err)

	a.True(
		len(translations) > 0,
		"Translations must contain at least one locale",
	)
}

func TestBuildLocales_DefaultLocale(t *testing.T) {
	a := assert.New(t)

	translations, _ := locales.BuildLocales()

	haveDefaultLocale := false
	for _, localeTranslations := range translations {
		if localeTranslations.Default {
			haveDefaultLocale = true
			break
		}
	}

	a.True(
		haveDefaultLocale,
		`Translations must contain "%s" locale`,
		locales.DefaultLocale,
	)
}

func translationKeys() []string {
	value := reflect.ValueOf(locales.TranslationKeys)
	keys := make([]string, 0, value.NumField())
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		keys = append(keys, field.String())
	}
	return keys
}

func defaultLocaleTranslations() models.LocaleTranslations {
	translations, err := locales.BuildLocales()
	if nil != err {
		panic(err)
	}

	for _, localeTranslations := range translations {
		if localeTranslations.Default {
			return localeTranslations
		}
	}
	panic("Default locale not found")
}

func TestBuildLocales_AllKeys(t *testing.T) {
	localeTranslations := defaultLocaleTranslations()

	for _, translationKey := range translationKeys() {
		translationExists := false
		for _, translation := range localeTranslations.Translations {
			if translation.Key == translationKey {
				translationExists = true
			}
		}
		if !translationExists {
			t.Errorf("Translation key %s not found", translationKey)
		}
	}
}
