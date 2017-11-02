package templates

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type list []template.Template

// ParseDir walks the provided `dir`
// returns a list of found templates
func ParseDir(dir string) ([]template.Template, error) {
	var (
		list []template.Template
	)

	absoluteDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	if err := filepath.Walk(absoluteDir, func(path string, info os.FileInfo, err error) error {

		log.Print(path)
		fileEnding := path[strings.LastIndex(path, ".")+1:]
		if !info.IsDir() && fileEnding == "templ" {
			list, err = appendTo(list, path, info.Name())
			if err != nil {
				return err
			}
		}

		return err
	}); err != nil {
		return nil, err
	}

	return list, nil
}

func Find(name string, list []template.Template) (*template.Template, error) {
	for _, l := range list {
		if l.Name() == name {
			log.Println(l.Name())
			return &l, nil
		}
	}

	return nil, fmt.Errorf("Template '%v' could not be found.", name)
}

func appendTo(l list, path, name string) (new list, err error) {

	templStr, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	t, err := template.New(name).Parse(string(templStr))
	if err != nil {
		return
	}

	new = append(l, *t)

	return
}
