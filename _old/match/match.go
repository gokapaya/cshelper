package match

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gokapaya/cshelper/csv"
)

// Pair represents a matched pair of users
type Pair struct {
	Giftee csv.User
	Santa  csv.User
}

func (p Pair) matchReport() string {
	return fmt.Sprintf("----\n%v -> %v\n\tinternational? %v -> %v\n\tcountry %v -> %v\n", p.Santa.Username, p.Giftee.Username, p.Santa.International, p.Giftee.International, p.Santa.Address.Country, p.Giftee.Address.Country)
}

// MatchUsers is the main function of this package
func MatchUsers(ul []csv.User) ([]Pair, error) {
	var (
		// lenL  = len(ul)
		allSantas       = make([]csv.User, len(ul))
		allGiftees      = make([]csv.User, len(ul))
		nationalSanta   = make(map[string][]csv.User)
		gifteeByCountry = make(map[string][]csv.User)
		countriesFound  = make(map[string]int)
		pairs           []Pair
	)

	copy(allSantas, ul)
	copy(allGiftees, ul)

	// sort national santas and giftees by country
	for _, u := range ul {
		if !u.International {
			nationalSanta[u.Address.Country] = append(nationalSanta[u.Address.Country], u)
		}
		gifteeByCountry[u.Address.Country] = append(gifteeByCountry[u.Address.Country], u)
		countriesFound[u.Address.Country]++
	}

	// match national santas
	for k := range countriesFound {
		santas := nationalSanta[k]

		if santas != nil {
			validGiftees := gifteeByCountry[k]
			for _, santa := range santas {
				// log.Printf("matching for santa: %v...\n", santa.Username)
				p := Pair{Santa: santa, Giftee: csv.User{Username: "nil"}}

				for !isValidPair(p) && !allSameSub(validGiftees, santa.RepSubreddit) {
					p.randomMatch(validGiftees)
				}
				if p.Giftee.Username != "nil" {
					allSantas = remove(allSantas, p.Santa)
					allGiftees = remove(allGiftees, p.Giftee)
					validGiftees = remove(validGiftees, p.Giftee)
					pairs = append(pairs, p)
				}
			}
		}
	}

	// match remaining santas
	for _, santa := range allSantas {
		// log.Printf("matching for santa: %v...\n", santa.Username)
		p := Pair{Santa: santa, Giftee: csv.User{Username: "nil"}}
		for !isValidPair(p) && !allSameSub(allGiftees, santa.RepSubreddit) {
			p.randomMatch(allGiftees)
		}
		if p.Giftee.Username != "nil" {
			allSantas = remove(allSantas, p.Santa)
			allGiftees = remove(allGiftees, p.Giftee)
			pairs = append(pairs, p)
		}
	}

	return pairs, nil
}

func (p *Pair) randomMatch(ul []csv.User) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	max := len(ul) - 1
	if max > 1 {
		p.Giftee = ul[r.Intn(max)]
		return
	}
	p.Giftee = ul[max]
}

func isValidPair(p Pair) bool {
	// TODO: is that sufficient?
	if p.Giftee.Username != "nil" {
		if p.Giftee.RepSubreddit != p.Santa.RepSubreddit {
			return true
		}
		return false
	}
	return false
}

func allSameSub(ul []csv.User, sub string) bool {
	for _, u := range ul {
		if u.RepSubreddit != sub {
			return false
		}
	}
	return true
}

// TODO: https://github.com/golang/go/wiki/SliceTricks
func remove(ul []csv.User, u csv.User) (new []csv.User) {
	for _, e := range ul {
		if e.Username != u.Username {
			new = append(new, e)
		}
	}
	return
}

func okay(p Pair) bool {
	if !p.Santa.International && !sameCountry(p) {
		return false
	}
	return true
}

func sameCountry(p Pair) bool {
	if p.Santa.Address.Country == p.Giftee.Address.Country {
		return true
	}
	return false
}
