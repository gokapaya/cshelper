package tmpl

import (
	"html/template"
	"path/filepath"

	"github.com/inconshreveable/log15"
)

var (
	Log  = log15.New()
	tmpl = template.New("global")
)

func Load() {
	var (
		rootDir = "./doc"
	)
	Log.Debug("walking tree for templates", "root", rootDir)

	docDir, err := filepath.Abs(rootDir)
	if err != nil {
		Log.Error("unable to get absolute doc path", "err", err)
		return
	}

	tmpl, err = template.ParseGlob(filepath.Join(docDir, "*.tmpl"))
	if err != nil {
		Log.Error("unable to parse template files")
		return
	}
}

func List() {
	Log.Info("listing registered templates")
	for _, t := range tmpl.Templates() {
		println(t.Name())
	}
}

func Lookup(name string) *template.Template {
	Log.Warn("this is a naive implementation")
	return tmpl.Lookup(name)
}
