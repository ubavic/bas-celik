package translation

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/ubavic/bas-celik/v2/localization"
)

var translations map[localization.Language]map[string]string

var currentLanguage localization.Language

func SetTranslations(embedFS embed.FS) error {
	translations = make(map[localization.Language]map[string]string)

	languages := []localization.Language{localization.SrLatin, localization.SrCyrillic, localization.En}
	for _, lang := range languages {
		langJson, err := embedFS.ReadFile("embed/translation/" + string(lang) + ".json")
		if err != nil {
			return fmt.Errorf("reading %s translation: %w", lang, err)
		}

		langMap := make(map[string]string)

		err = json.Unmarshal(langJson, &langMap)
		if err != nil {
			return fmt.Errorf("%s translation unmarshal: %w", lang, err)
		}

		translations[lang] = langMap
	}

	return nil
}

func SetLanguage(lang int) {
	if lang == 2 {
		currentLanguage = localization.En
	} else if lang == 1 {
		currentLanguage = localization.SrCyrillic
	} else {
		currentLanguage = localization.SrLatin
	}
}

func CurrentLanguage() localization.Language {
	return currentLanguage
}

func Translate(id string, vals ...any) string {
	translation := translations[currentLanguage][id]

	if len(vals) == 0 {
		return translation
	} else {
		return fmt.Sprintf(translation, vals...)
	}
}

func EnglishTranslation(id string) string {
	return translations[localization.En][id]
}
