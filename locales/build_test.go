package locales_test

import (
	"testing"

	"github.com/c0va23/redirector/locales"

	"github.com/stretchr/testify/assert"
)

const defaultLocale = "en"

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
		if defaultLocale == localeTranslations.Code {
			haveDefaultLocale = true
			break
		}
	}

	a.True(
		haveDefaultLocale,
		`Translations must contain "%s" locale`,
		defaultLocale,
	)
}
