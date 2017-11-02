package csv

import "github.com/gokapaya/cshelper/errlist"

// GetPairs reads a SantaList.csv file and returns a map of matches
func GetPairs(path string) (map[string]string, error) {
	var (
		el    errlist.ErrList
		pairs = make(map[string]string)
	)

	el = parseCSVFile(path, func(rec []string) error {
		pairs[rec[1]] = rec[0]
		return nil
	})

	if el.NotEmpty() {
		return pairs, el
	}
	return pairs, nil
}
