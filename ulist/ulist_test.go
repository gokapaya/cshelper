package ulist

import (
	"os"
	"testing"

	"github.com/naoina/toml"
	"github.com/stretchr/testify/assert"
)

var ul Ulist

func setup() {
	fd, err := os.Open("ulist_test.toml")
	if err != nil {
		panic(err)
	}
	dec := toml.NewDecoder(fd)
	var tomlData tomlUlist
	if err := dec.Decode(&tomlData); err != nil {
		panic(err)
	}
	ul = *NewUlist(tomlData.Users)
}

func init() {
	setup()
}

func TestGetByName(t *testing.T) {
	var exp = &User{
		Username: "test1", ShareName: false, AgreeRules: true,
		RepSubreddit: "/r/AraragiGirls", MessageToSanta: "test1 test",
		Watchlist: "",
		Address: Address{
			Fullname: "Jon Doe",
			Street1:  "Doestr. 1", Street2: "",
			City: "Doe City", State: "DX",
			Country: "United States", Zipcode: "00000",
		},
		Regift: true, International: false, match: "",
	}

	result := ul.GetByName("test1")
	assert.Equal(t, exp, result)
}

func TestFilter(t *testing.T) {
	var exp = *NewUlist([]User{
		*ul.GetByName("test1"),
	})

	result := ul.Filter(func(u User) bool {
		return CompareUsernames(u.Username, "test1")
	})

	assert.Equal(t, exp, result)
}

func TestGetByCountry(t *testing.T) {
	var exp = *NewUlist([]User{
		*ul.GetByName("test1"),
	})

	result := ul.GetByCountry("US")

	assert.Equal(t, exp, result)
}
