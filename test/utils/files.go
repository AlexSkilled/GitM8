package utils

import (
	"io/ioutil"
	"strings"

	"github.com/sirupsen/logrus"
)

func GetFiles(dir string) []string {
	filesMigrations, err := ioutil.ReadDir(dir)
	if err != nil {
		logrus.Fatal(err)
	}
	files := make([]string, 0, len(filesMigrations))
	for _, f := range filesMigrations {
		name := f.Name()
		if strings.Contains(name, ".http") {
			files = append(files, name)
		}
	}
	for i, item := range files {
		file, err := ioutil.ReadFile(dir + item)
		if err != nil {

		}
		files[i] = string(file)
	}

	return files
}
