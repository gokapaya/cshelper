package match

import (
	"encoding/csv"
	"os"
)

func WriteToCSV(dry bool, pairs []Pair, filename string) error {
	w := csv.NewWriter(os.Stdout)
	if !dry {
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		w = csv.NewWriter(f)
	}

	w.Write([]string{"User", "User's Closet Santa"})
	for _, p := range pairs {
		if err := w.Write([]string{p.Giftee.Username, p.Santa.Username}); err != nil {
			return err
		}
	}
	w.Flush()

	return w.Error()
}
