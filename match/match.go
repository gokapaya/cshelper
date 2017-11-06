package match

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/clyphub/munkres"
	"github.com/gocarina/gocsv"
	"github.com/gokapaya/cshelper/ulist"
	"github.com/inconshreveable/log15"
	"github.com/pkg/errors"
)

var Log = log15.New()

type Pair struct {
	Santa  *ulist.User
	Giftee *ulist.User
}

type data struct {
	Giftee string `csv:"user"`
	Santa  string `csv:"santa"`
}

// Match takes a list of Users and returns Pairs of Santa and Giftee
func Match(ul *ulist.Ulist) ([]Pair, error) {
	Log.Info("matching users")
	if ul.Len() < 2 {
		return nil, errors.New("user list has < 2 entries")
	}

	m := munkres.NewMatrix(ul.Len())
	m.A = costMatrix(ul, func(santa, giftee ulist.User) int64 {
		var cost int64 = 0
		if santa.Username == giftee.Username {
			cost += 100
		}

		if !santa.International {
			if !ulist.SameCountry(santa.Address.Country, giftee.Address.Country) {
				if !ulist.SameRegion(santa.Address.Country, giftee.Address.Country) {
					cost += 100
				} else {
					cost += 20
				}
			}
		}
		return cost
	})

	Log.Debug("running munkres")
	result := munkres.ComputeMunkresMin(m)
	// m.Print()
	// printRowCol(ul, result)

	var pairs = make([]Pair, 0)
	for _, rowcol := range result {
		pairs = append(pairs, Pair{
			Santa:  ul.Get(rowcol.Row),
			Giftee: ul.Get(rowcol.Col),
		})
	}
	return pairs, nil
}

func costMatrix(ul *ulist.Ulist, costFn func(ulist.User, ulist.User) int64) []int64 {
	Log.Debug("calculating cost matrix", "total_users", ul.Len())
	var m = make([][]int64, ul.Len())

	ul.Iter(func(i int, a ulist.User) error {
		m[i] = make([]int64, ul.Len())

		ul.Iter(func(k int, b ulist.User) error {
			var cost int64 = 0

			cost += costFn(a, b)
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			cost += int64(r.Int63n(20))

			m[i][k] = cost
			return nil

		})
		return nil
	})

	return matrixToSlice(m)
}

func matrixToSlice(matrix [][]int64) []int64 {
	var slice []int64
	for _, array := range matrix {
		for _, num := range array {
			slice = append(slice, num)
		}
	}
	return slice
}

// Eval takes a []Pair and evaluates it agains a set of rules.
func Eval(pairings []Pair) error {
	Log.Debug("running pair evalutation")
	for _, p := range pairings {
		// check not same name

		if p.Santa.Username == p.Giftee.Username {
			Log.Warn("evaluation failed")
			return errors.Errorf("same person\n\n%s == %s", p.Santa.Username, p.Giftee.Username)
		}

		if !p.Santa.International {
			if !ulist.SameRegion(
				p.Santa.Address.Country,
				p.Giftee.Address.Country,
			) {
				Log.Warn("evaluation failed")
				return errors.Errorf("santa doesn't want international but has to send out of his region\n\nSanta's country: %s\nGiftee's country: %s", p.Santa.Address.Country, p.Giftee.Address.Country)
			}
		}
	}
	Log.Info("evaluation successful")
	return nil
}

func printRowCol(ul *ulist.Ulist, result []munkres.RowCol) {
	println("result:")
	for _, rc := range result {
		fmt.Printf("% 3v <> % 3v\n\n",
			ul.Get(rc.Row).Username,
			ul.Get(rc.Col).Username,
		)
	}
}

func SavePairings(fpath string, pairings []Pair) error {
	defer Log.Info("pair csv saved", "file", fpath)
	// save .csv file
	fd, err := os.OpenFile(fpath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "unable to open %q for writing", fpath)
	}
	var d []*data
	for _, p := range pairings {
		d = append(d, &data{Santa: p.Santa.Username, Giftee: p.Giftee.Username})
	}
	return gocsv.MarshalFile(d, fd)
}

func LoadPairings(fpath string, ul ulist.Ulist) ([]Pair, error) {
	defer Log.Info("pair csv loaded", "file", fpath)
	fd, err := os.Open(fpath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to open %q for reading", fpath)
	}

	var p []data
	if err := gocsv.UnmarshalFile(fd, p); err != nil {
		return nil, err
	}

	var pairs []Pair
	for _, datap := range p {
		pair := Pair{
			Santa:  ul.GetByName(datap.Santa),
			Giftee: ul.GetByName(datap.Giftee),
		}
		if pair.Santa == nil {
			return nil, errors.Errorf("username not found %q", datap.Santa)
		}
		if pair.Giftee != nil {
			return nil, errors.Errorf("username not found %q", datap.Giftee)
		}
		pairs = append(pairs, pair)
	}
	return pairs, nil
}
