package ulist

import (
	"os"
	"strings"

	"github.com/pirsquare/country-mapper"
)

var infoClient *country_mapper.CountryInfoClient

func init() {
	cl, err := country_mapper.Load()
	if err != nil {
		Log.Crit("unable to load country_mapper data", "err", err)
		os.Exit(1)
	}
	infoClient = cl
}

func (ul *Ulist) GetByCountry(country string) *Ulist {
	return ul.Filter(func(u User) bool {
		return SameCountry(u.Address.Country, country)
	})
}

const (
	countryUS = "United States"
	countryUK = "United Kingdom"
	countryDE = "Deutschland"
	countryNL = "Netherlands"
)

var countryReplacer = strings.NewReplacer(
	"United States", countryUS,
	"Unites States", countryUS,
	"USA", countryUS,
	"US", countryUS,
	"United States of America", countryUS,
	"North America", countryUS,

	"United Kingdom", countryUK,
	"Great Britain", countryUK,
	"UK", countryUK,

	"Deutschland", countryDE,
	"Germany", countryDE,

	"The Netherlands", countryNL,
)

func NormalizedCountry(country string) string {
	return countryReplacer.Replace(country)
}

func ClusterCountry(country string) string {
	cleaned := strings.TrimSpace(country)

	d := infoClient.MapByName(cleaned)
	if d == nil {
		Log.Error("country not found", "str", country)
		return ""
	}
	return d.Region
}

func SameCountry(a, b string) bool {
	return NormalizedCountry(a) == NormalizedCountry(b)
}

func SameRegion(a, b string) bool {
	clA := ClusterCountry(NormalizedCountry(a))
	clB := ClusterCountry(NormalizedCountry(b))
	return (clA == clB)
}
