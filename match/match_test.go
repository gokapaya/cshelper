package match

import (
	"testing"

	"github.com/gokapaya/cshelper/ulist"
	"github.com/stretchr/testify/assert"
)

var ul = ulist.NewUlist([]ulist.User{
	ulist.User{
		Username:      "t1",
		International: true, Regift: true,
		Address:      ulist.Address{Country: "US"},
		RepSubreddit: "Chuunibyou",
	},
	ulist.User{
		Username:      "t2",
		International: false, Regift: false,
		Address:      ulist.Address{Country: "US"},
		RepSubreddit: "aria",
	},
	ulist.User{
		Username:      "t3",
		International: true, Regift: true,
		Address:      ulist.Address{Country: "Germany"},
		RepSubreddit: "k_on",
	},
	ulist.User{
		Username:      "t4",
		International: true, Regift: true,
		Address:      ulist.Address{Country: "United States"},
		RepSubreddit: "Evangelion",
	},
	ulist.User{
		Username:      "t5",
		International: false, Regift: true,
		Address:      ulist.Address{Country: "United Kingdom"},
		RepSubreddit: "chuunibyou",
	},
})

func TestMatching(t *testing.T) {
	pairings, err := Match(ul)

	assert.Nil(t, err)
	assert.NotEmpty(t, pairings)

	assert.Condition(t, func() bool {
		if err := Eval(pairings); err != nil {
			t.Log(err)
			return false
		}
		return true
	})
}
