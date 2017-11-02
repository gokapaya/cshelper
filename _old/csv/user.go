package csv

import (
	"fmt"
	"net/url"
)

type Address struct {
	Fullname string
	Street1  string
	Street2  string
	City     string
	State    string
	Country  string
	Zipcode  string
}

func (a *Address) String() string {
	return fmt.Sprintf("%v\n%v\n%v\n%v %v\n%v %v", a.Fullname, a.Street1, a.Street2, a.Zipcode, a.City, a.State, a.Country)
}

func newAddress(rec []string) Address {
	return Address{
		Fullname: rec[0],
		Street1:  rec[1],
		Street2:  rec[2],
		City:     rec[3],
		State:    rec[4],
		Country:  rec[5],
		Zipcode:  rec[6],
	}
}

// User represents a participant in the ClosetSanta
type User struct {
	Username       string
	ShareName      bool
	AgreeRules     bool
	RepSubreddit   string
	MessageToSanta string
	Watchlist      *url.URL
	Address        Address
	Regift         bool
	International  bool
}

func parseBool(yesno string) bool {
	if yesno == "Yes" {
		return true
	}
	return false
}

func parseURL(urlStr string) *url.URL {
	url, err := url.Parse(urlStr)
	if err != nil {
		return nil
	}

	return url
}

// TODO: errors?
func newUser(rec []string) (*User, error) {
	return &User{
		Username:       rec[1],
		ShareName:      parseBool(rec[2]),
		AgreeRules:     parseBool(rec[3]),
		RepSubreddit:   rec[4],
		MessageToSanta: rec[5],
		Watchlist:      parseURL(rec[6]),
		Address:        newAddress(rec[7:]),
		Regift:         parseBool(rec[14]),
		International:  parseBool(rec[15]),
	}, nil
}

// Key implements the boltclient.Thing interface
func (u User) Key() []byte {
	return []byte(u.Username)
}

// Bucket implements the boltclient.Thing interface
func (u User) Bucket() []byte {
	return []byte("giftees")
}
