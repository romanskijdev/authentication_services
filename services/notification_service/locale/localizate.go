package locale

import (
	_ "embed"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

//go:embed locale.en.toml
var localeEn string

//go:embed locale.ru.toml
var localeRu string

var LanguageMapTranslate = map[string]bool{
	"en": true,
	"ru": true,
}

func I8nInit() *i18n.Bundle {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseMessageFileBytes([]byte(localeEn), "locale.en.toml")
	bundle.MustParseMessageFileBytes([]byte(localeRu), "locale.ru.toml")
	return bundle
}

func LocaleConvert(langCode *string, tag string, bundle *i18n.Bundle) (string, error) {
	langCodeStr := "en"
	if langCode != nil {
		langCodeStr = *langCode
		if !LanguageMapTranslate[langCodeStr] {
			langCodeStr = "en"
		}
	}
	localizer := i18n.NewLocalizer(bundle, langCodeStr)
	translation, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: tag})
	if err != nil {
		logrus.Errorln("🔴 error localize: ", err)
		// В случае ошибки возвращаем неизмененную подстроку
		return "", err
	}
	return translation, nil
}
