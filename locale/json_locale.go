package locale

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

var languagePackRegex = regexp.MustCompile("^(bundle_).+\\.(json)$")

// NewFromDir create a Locale instance from files with name format of `bundle_[LANGUAGE_TAG].json` inside the directory.
func NewFromDir(dirname string) (Locale, error) {
	localeBundle := jsonLocaleBundle{
		bundle: make(map[Language]map[string]string),
	}

	err := filepath.Walk(dirname, func(filePath string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		fileName := info.Name()
		if languagePackRegex.MatchString(fileName) {
			languagePackFile, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer languagePackFile.Close()

			languageTag := Language(fileName[7:9])
			languagePackBytes, err := ioutil.ReadAll(languagePackFile)
			if err != nil {
				return fmt.Errorf("error on reading locale bundle file: %w", err)
			}
			localPackMap := make(map[string]string)
			err = json.Unmarshal(languagePackBytes, &localPackMap)
			if err != nil {
				return fmt.Errorf("error on parsing locale bundle file: %w", err)
			}

			localeBundle.bundle[languageTag] = localPackMap
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &localeBundle, nil
}

type jsonLocaleBundle struct {
	bundle map[Language]map[string]string
}

func (j *jsonLocaleBundle) Message(lang Language, key string, params ...interface{}) string {
	// guard with default value
	if lang == "" || len(lang) != 2 {
		lang = English
	}

	languagePack := j.bundle[lang]
	if languagePack == nil {
		lang = English
		languagePack = j.bundle[lang]

		if languagePack == nil {
			return "locale file problems!!!"
		}
	}

	// do not use fmt for simple messages
	if len(params) == 0 {
		return languagePack[key]
	}

	return fmt.Sprintf(languagePack[key], params...)
}
