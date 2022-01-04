package langs

import (
	"context"
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"

	"gitlab-tg-bot/utils"

	"github.com/hashicorp/go-uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var AvailableLangs []string

// key - message header; value - map of translations language -> message
var languagesHandler map[string]map[string]string

var dL string

func GetDefaultLocale() string {
	return dL
}

// Init - creates an instance of redis connection
func Init(defaultLanguage string, langPaths []string) error {
	dL = defaultLanguage

	dir, _ := os.Getwd()
	dir += "/internal/message-handling/"

	// TODO add langPaths to files
	//

	languages, err := ioutil.ReadDir(dir)
	if err != nil {
		logrus.Fatal(err)
	}

	var defaultLanguageFile fs.FileInfo

	files := make([]fs.FileInfo, 0, len(languages))

	for _, f := range languages {
		name := f.Name()
		if strings.HasSuffix(name, ".json") {
			if strings.HasPrefix(name, dL) {
				defaultLanguageFile = f
			} else {
				files = append(files, f)
			}
		}
	}

	if defaultLanguageFile == nil {
		panic("No default language has been set")
	} else {
		languagesHandler = make(map[string]map[string]string)
		err = addLanguage(dir+defaultLanguageFile.Name(), dL)
		if err != nil {
			return err
		}
	}

	for _, item := range files {
		locale := item.Name()
		localeName := locale[0:strings.Index(locale, ".")]
		if err = addLanguage(dir+locale, localeName); err != nil {
			logrus.Error(err, "Ошибка при загрузки локали")
		}
	}
	return nil
}

func addLanguage(pathToFile string, language string) error {
	AvailableLangs = append(AvailableLangs, language)
	content, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		return errors.Wrap(err, "Error while reading default language message markups")
	}
	var defaultMarkupRaw map[string]interface{}
	err = json.Unmarshal(content, &defaultMarkupRaw)
	if err != nil {
		return errors.Wrap(err, "Error while unmarshalling default markup")
	}
	markup := make(map[string]string)

	extractFullKeys("", defaultMarkupRaw, markup)

	for key, value := range markup {
		if languagesHandler[key[1:]] == nil {
			languagesHandler[key[1:]] = make(map[string]string)
		}
		languagesHandler[key[1:]][language] = value
	}
	return nil
}

func extractFullKeys(currentKey string, source map[string]interface{}, target map[string]string) {
	for key, value := range source {
		switch v := value.(type) {
		case map[string]interface{}:
			extractFullKeys(currentKey+":"+key, v, target)
		case string:
			target[currentKey+":"+key] = v
		}
	}
}

func Get(ctx context.Context, key string) (out string) {
	locale, err := utils.ExtractLocale(ctx)
	if err != nil {
		locale = dL
	}
	translations, ok := languagesHandler[key]
	if !ok {
		uid, _ := uuid.GenerateUUID()
		logrus.Warnf("Attention! Message key %s doesn't exist. Return exception key %s", key, uid)
		return uid
	}
	out, ok = translations[locale]
	if !ok {
		logrus.Warnf("Attention! For key %s Locale %s doesn't exist. Extracting from default", key, locale)
		out, ok = translations[dL]
		if !ok {
			uid, _ := uuid.GenerateUUID()
			logrus.Errorf("Error when tried to extract message with default locale: %s. Return exception key: %s", dL, uid)
			return uid
		}
	}

	return out
}

func GetWithLocale(locale, key string) (out string) {
	translations, ok := languagesHandler[key]
	if !ok {
		uid, _ := uuid.GenerateUUID()
		logrus.Warnf("Attention! Message key %s doesn't exist. Return exception key %s", key, uid)
		return uid
	}
	out, ok = translations[locale]
	if !ok {
		logrus.Warnf("Attention! For key %s Locale %s doesn't exist. Extracting from default", key, locale)
		out, ok = translations[dL]
		if !ok {
			uid, _ := uuid.GenerateUUID()
			logrus.Errorf("Error when tried to extract message with default locale: %s. Return exception key: %s", dL, uid)
			return uid
		}
	}

	return out
}
