package ulist

import (
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gocarina/gocsv"
	"github.com/pkg/errors"
)

type data struct {
	Timestamp      string `csv:"Timestamp"`
	Username       string `csv:"What is your reddit username?"`
	ShareUsername  string `csv:"Share Reddit Username with Santa"`
	AgreeRules     string `csv:"Do you understand and agree to the rules of this gift exchange?"`
	Subreddit      string `csv:"Subreddit to represent:"`
	MessageToSanta string `csv:"Message to your Santa"`
	Watchlist      string `csv:"Watchlist:"`
	Fullname       string `csv:"Name:"`
	Street1        string `csv:"Address:"`
	Street2        string `csv:"Address 2:"`
	City           string `csv:"City:"`
	State          string `csv:"State/Province:"`
	Country        string `csv:"Country:"`
	Zipcode        string `csv:"Post/Zip Code:"`
	Rematcher      string `csv:"Are you willing to be put down as a possible rematcher?"`
	ShipAbroad     string `csv:"Are you willing to ship abroad?"`
}

func (d *data) toUser() (*User, error) {
	var err error
	u := &User{
		Username:       d.Username,
		ShareName:      parseBool(d.ShareUsername),
		AgreeRules:     parseBool(d.AgreeRules),
		RepSubreddit:   d.Subreddit,
		MessageToSanta: d.MessageToSanta,
		Address: Address{
			Fullname: d.Fullname,
			Street1:  d.Street1,
			Street2:  d.Street2,
			City:     d.City,
			State:    d.State,
			Country:  d.Country,
			Zipcode:  d.Zipcode,
		},
		Regift:        parseBool(d.Rematcher),
		International: parseBool(d.ShipAbroad),
	}
	u.Watchlist, err = parseURL(d.Watchlist)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse watchlist url")
	}
	return u, nil
}

func ParseFile(fpath string) ([]User, error) {
	Log.Debug("parsing csv", "path", fpath)
	var ulist []User

	absFile, err := filepath.Abs(fpath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get absolute filepath")
	}

	fd, err := os.Open(absFile)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to open %q", absFile)
	}

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		return gocsv.LazyCSVReader(in)
	})

	var datas = []*data{}
	if err := gocsv.UnmarshalFile(fd, &datas); err != nil {
		return nil, errors.Wrap(err, "unable to read csv")
	}

	for _, d := range datas {
		u, err := d.toUser()
		if err != nil {
			return nil, errors.Wrap(err, "unable to convert entry to User")
		}
		ulist = append(ulist, *u)
	}

	return ulist, nil
}

func parseBool(yesno string) bool {
	if yesno == "Yes" {
		return true
	}
	return false
}

func parseURL(urlStr string) (string, error) {
	urlObj, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	return urlObj.String(), nil
}
