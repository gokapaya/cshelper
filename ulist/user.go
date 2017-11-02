package ulist

import (
	"fmt"
)

// User represents a participant in the ClosetSanta
type User struct {
	Username       string
	ShareName      bool
	AgreeRules     bool
	RepSubreddit   string
	MessageToSanta string
	Watchlist      string
	Address        Address
	Regift         bool
	International  bool
	match          string
}

func (u *User) GetMatch(ulist []User) *User {
	for _, user := range ulist {
		if u.match == user.Username {
			return &user
		}
	}
	return nil
}

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
