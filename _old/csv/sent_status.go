package csv

import "github.com/gokapaya/cshelper/errlist"

// ShippingStatusMap maps a username to a status
type ShippingStatusMap map[string]string

// ShippingStatus ...
type ShippingStatus struct {
	User                        string
	Sent, Fail, Resolved, Check bool
}

// GetShippingStatusInf returns a ShippingStatusMap from a CSV
func GetShippingStatusInf(path string) (ShippingStatusMap, error) {
	var (
		stat = make(ShippingStatusMap)
		el   errlist.ErrList
	)

	el = parseCSVFile(path, func(rec []string) error {
		stat[rec[0]] = rec[1]
		return nil
	})

	if el.NotEmpty() {
		return stat, el
	}
	return stat, nil
}
