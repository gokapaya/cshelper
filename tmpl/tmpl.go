package tmpl

import (
	"html/template"
	"path/filepath"

	"github.com/inconshreveable/log15"
)

var (
	Log = log15.New()
)

func Lookup(name string) *template.Template {
	path, err := filepath.Abs(name)
	if err != nil {
		Log.Error("unable to get absolute doc path", "err", err)
		return nil
	}

	t, err := template.ParseFiles(path)
	if err != nil {
		Log.Error("error parsing template", "err", err)
		return nil
	}
	return t
}
