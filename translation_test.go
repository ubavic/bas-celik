package main

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"maps"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
	"unicode"

	"github.com/ubavic/bas-celik/v2/localization"
)

var placeholderPattern = regexp.MustCompile(`%[a-zA-Z]`)
var translationKeyPattern = regexp.MustCompile(`^[a-z][a-zA-Z0-9]*(\.[a-zA-Z0-9_-]+)+$`)
var fileNamePattern = regexp.MustCompile(`\.(go|json|md|txt|xml|pdf|pem|crt|cer|xlsx|png|jpg|ttf|exe|dll|so|dylib|yml|yaml|toml)$`)
var latinInCyrillicPattern = regexp.MustCompile(`PCSC-lite|PKCS#11|Smartbox|Excel|PDF|XML|PIN|ATR|web`)

func loadTranslations(t *testing.T) map[localization.Language]map[string]string {
	t.Helper()

	all := make(map[localization.Language]map[string]string, len(localization.Languages))

	for _, lang := range localization.Languages {
		data, err := embedFS.ReadFile("embed/translation/" + string(lang) + ".json")
		if err != nil {
			t.Fatalf("reading %s translation: %v", lang, err)
		}

		translations := make(map[string]string)
		if err := json.Unmarshal(data, &translations); err != nil {
			t.Fatalf("unmarshalling %s translation: %v", lang, err)
		}

		all[lang] = translations
	}

	return all
}

func Test_TranslationsDefineSameKeys(t *testing.T) {
	all := loadTranslations(t)

	for _, lang := range localization.Languages {
		if lang == localization.En {
			continue
		}

		for _, key := range slices.Sorted(maps.Keys(all[localization.En])) {
			if _, ok := all[lang][key]; !ok {
				t.Errorf("%s does not define key %q", lang, key)
			}
		}

		for _, key := range slices.Sorted(maps.Keys(all[lang])) {
			if _, ok := all[localization.En][key]; !ok {
				t.Errorf("%s defines key %q that %s does not", lang, key, localization.En)
			}
		}
	}
}

func Test_TranslationKeysAreWellFormed(t *testing.T) {
	all := loadTranslations(t)

	for _, lang := range localization.Languages {
		for _, key := range slices.Sorted(maps.Keys(all[lang])) {
			if !translationKeyPattern.MatchString(key) {
				t.Errorf("%s key %q is not of the form group.name", lang, key)
			}

			if fileNamePattern.MatchString(key) {
				t.Errorf("%s key %q looks like a file name", lang, key)
			}
		}
	}
}

func Test_TranslationsHaveValuesAndMatchingPlaceholders(t *testing.T) {
	all := loadTranslations(t)

	for _, key := range slices.Sorted(maps.Keys(all[localization.En])) {
		expected := placeholderPattern.FindAllString(all[localization.En][key], -1)

		for _, lang := range localization.Languages {
			value, ok := all[lang][key]
			if !ok {
				continue // Reported by Test_TranslationsDefineSameKeys.
			}

			if strings.TrimSpace(value) == "" {
				t.Errorf("%s has an empty value for key %q", lang, key)
				continue
			}

			got := placeholderPattern.FindAllString(value, -1)
			if !slices.Equal(got, expected) {
				t.Errorf("%s key %q has placeholders %v, but %s has %v", lang, key, got, localization.En, expected)
			}
		}
	}
}

func Test_TranslationKeysUsedInCodeAreDefined(t *testing.T) {
	english := loadTranslations(t)[localization.En]

	groups := make(map[string]bool)
	for key := range english {
		group, _, _ := strings.Cut(key, ".")
		groups[group] = true
	}

	for _, literal := range goStringLiterals(t) {
		if !translationKeyPattern.MatchString(literal) || fileNamePattern.MatchString(literal) {
			continue
		}

		group, _, _ := strings.Cut(literal, ".")
		if !groups[group] {
			continue
		}

		if _, ok := english[literal]; !ok {
			t.Errorf("code references translation key %q, but no translation file defines it", literal)
		}
	}
}

func Test_SerbianTranslationsDoNotMixScripts(t *testing.T) {
	all := loadTranslations(t)

	for _, key := range slices.Sorted(maps.Keys(all[localization.SrLatin])) {
		value := all[localization.SrLatin][key]

		for _, char := range value {
			if unicode.Is(unicode.Cyrillic, char) {
				t.Errorf("%s key %q contains Cyrillic %q (%U): %q", localization.SrLatin, key, char, char, value)
				break
			}
		}
	}

	for _, key := range slices.Sorted(maps.Keys(all[localization.SrCyrillic])) {
		value := all[localization.SrCyrillic][key]

		stripped := placeholderPattern.ReplaceAllString(value, "")
		stripped = latinInCyrillicPattern.ReplaceAllString(stripped, "")

		for _, char := range stripped {
			if unicode.Is(unicode.Latin, char) {
				t.Errorf("%s key %q contains Latin %q (%U): %q", localization.SrCyrillic, key, char, char, value)
				break
			}
		}
	}
}

func goStringLiterals(t *testing.T) []string {
	t.Helper()

	var literals []string
	fileSet := token.NewFileSet()

	err := filepath.WalkDir(".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			if name := entry.Name(); name == ".git" || name == "vendor" {
				return fs.SkipDir
			}

			return nil
		}

		if filepath.Ext(path) != ".go" {
			return nil
		}

		parsed, err := parser.ParseFile(fileSet, path, nil, 0)
		if err != nil {
			return err
		}

		ast.Inspect(parsed, func(node ast.Node) bool {
			literal, ok := node.(*ast.BasicLit)
			if !ok || literal.Kind != token.STRING {
				return true
			}

			if value, err := strconv.Unquote(literal.Value); err == nil {
				literals = append(literals, value)
			}

			return true
		})

		return nil
	})

	if err != nil {
		t.Fatalf("scanning Go sources: %v", err)
	}

	return literals
}
