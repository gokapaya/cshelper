package ulist

import (
	"net/url"
	"os"
	"path/filepath"

	"github.com/gocarina/gocsv"
	"github.com/pkg/errors"
)

type data struct {
	Timestamp      string `csv:"timestamp"`
	Username       string `csv:"username"`
	ShareUsername  string `csv:"share_username"`
	AgreeRules     string `csv:"agree_rules"`
	Subreddit      string `csv:"subreddit"`
	MessageToSanta string `csv:"message_to_santa"`
	Watchlist      string `csv:"watchlist"`
	Fullname       string `csv:"fullname"`
	Street1        string `csv:"street_1"`
	Street2        string `csv:"street_2"`
	City           string `csv:"city"`
	State          string `csv:"state"`
	Country        string `csv:"country"`
	Zipcode        string `csv:"zip"`
	Rematcher      string `csv:"rematcher"`
	ShipAbroad     string `csv:"ship_abroad"`
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
