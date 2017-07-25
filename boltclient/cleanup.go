package boltclient

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/gokapaya/cshelper/csv"
	"github.com/gokapaya/cshelper/errlist"
)

func (c *Client) CleanDatabase() error {
	var (
		el      errlist.ErrList
		cleaned int
	)
	tx, err := c.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	b := tx.Bucket([]byte(UserBucket))
	if b == nil {
		return errlist.Err("Bucket not found")
	}

	cr := b.Cursor()
	for k, v := cr.First(); k != nil; k, v = cr.Next() {
		var u csv.User
		if err := json.Unmarshal(v, &u); err != nil {
			el = append(el, err)
			continue
		}

		// 		if fixCountry(&u) {
		// 			cleaned++
		// 			log.Println(u.Address.Country)
		// 			data, err := json.Marshal(u)
		// 			if err != nil {
		// 				el = append(el, err)
		// 				continue
		// 			}
		// 			v = data
		// 			if err := cr.Bucket().Put(k, v); err != nil {
		// 				el = append(el, err)
		// 				continue
		// 			}
		// 			log.Println("---")
		// 		}
		cleaned++
		log.Println(u.Username, "-->", strings.ToLower(u.Username))
		data, err := json.Marshal(u)
		if err != nil {
			el = append(el, err)
			continue
		}
		v = data
		if err := cr.Bucket().Put([]byte(strings.ToLower(string(k))), v); err != nil {
			el = append(el, err)
			continue
		}
		log.Println("---")
	}

	log.Printf("cleaned %v entries", cleaned)

	if el.NotEmpty() {
		return el
	}

	return tx.Commit()
}

type country string

func (c country) has(s string) bool {
	return strings.Contains(string(c), s)
}

func fixCountry(u *csv.User) bool {
	var c = country(u.Address.Country)

	if string(c) == "United States of America" {
		return false
	}

	if c.has("US") || c.has("United States") || c.has("usa") || c.has("Usa") || c.has("North America") {
		log.Println(c)
		u.Address.Country = "United States of America"
		return true
	}

	return false
}
