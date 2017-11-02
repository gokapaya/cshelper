package ulist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	var exp = []User{
		User{
			Username: "test1", ShareName: false, AgreeRules: true,
			RepSubreddit: "/r/AraragiGirls", MessageToSanta: "test1 test",
			Watchlist: "https://example.com",
			Address: Address{
				Fullname: "John Doe",
				Street1:  "jon doe blv", Street2: "",
				City: "Doetown", State: "DX",
				Country: "United States", Zipcode: "00000",
			},
			Regift: true, International: false, match: "",
		},
		User{
			Username: "test2", ShareName: false, AgreeRules: true,
			RepSubreddit: "/r/LoveLive", MessageToSanta: "test2 test",
			Watchlist: "http://example.com",
			Address: Address{
				Fullname: "Max Mustermann",
				Street1:  "Musterstr. 1", Street2: "Zimmer 1",
				City: "Neustadt", State: "Neuland",
				Country: "Germany", Zipcode: "99999",
			},
			Regift: false, International: false, match: "",
		},
	}

	result, err := ParseFile("parse_test.csv")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, result, exp)
}
