package ulist

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
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

func (u *User) String() string {
	// const fmtString = `
	// Username: {{ .Username }}
	// RepSubreddit: {{ .RepSubreddit }}
	// Watchlist: {{ .Watchlist }}
	// MessageToSanta:
	// ---
	// {{ .MessageToSanta }}
	// ---
	// AgreeRules: {{ printf "%t" .AgreeRules }}, ShareName: {{ printf "%t" .ShareName }}
	// Regift: {{ printf "%t" .Regift }}, International: {{ printf "%t" .International }}
	// Address:
	// ---
	// {{ .Address.String }}
	// ---
	// `
	const fmtString = `
Username: {{ .Username }}
  RepSubreddit: {{ .RepSubreddit }}
  Watchlist: {{ .Watchlist }}
  Address:
---
{{ .Address.String }}
---
  AgreeRules: {{ printf "%t" .AgreeRules }}
  ShareName: {{ printf "%t" .ShareName }}
  Regift: {{ printf "%t" .Regift }}
  International: {{ printf "%t" .International }}
  MessageToSanta:
--
{{ .MessageToSanta }}
---
`
	t := template.Must(template.New("fmt").Parse(fmtString))

	var buf bytes.Buffer
	if err := t.Execute(&buf, u); err != nil {
		panic(err)
	}
	return buf.String()
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

func CompareUsernames(u1, u2 string) bool {
	return strings.ToLower(u1) == strings.ToLower(u2)
}
