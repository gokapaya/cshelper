package csv

import (
	"encoding/csv"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/gokapaya/cshelper/errlist"
)

func parseCSVFile(path string, fn func([]string) error) errlist.ErrList {
	var el errlist.ErrList

	absoluteDir, err := filepath.Abs(path)
	if err != nil {
		el = append(el, err)
	}

	rawData, err := ioutil.ReadFile(absoluteDir)
	if err != nil {
		el = append(el, err)
	}

	r := csv.NewReader(strings.NewReader(string(rawData)))

	var i int
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			el = append(el, err)
		}
		if i > 0 {
			err := fn(rec)
			if err != nil {
				el = append(el, err)
			}
		}
		i++
	}

	return el
}

// GetUserList returns a list of users found in the CSV file
func GetUserList(path string) ([]User, error) {
	var (
		ul []User
		el errlist.ErrList
	)

	el = parseCSVFile(path, func(rec []string) error {
		u, err := newUser(rec)
		if err != nil {
			el = append(el, err)
		}
		ul = append(ul, *u)
		return nil
	})

	log.Printf("Found %v entries.", len(ul))

	if el.NotEmpty() {
		return ul, el
	}
	return ul, nil
}
