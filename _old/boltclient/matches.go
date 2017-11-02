package boltclient

import (
	"encoding/json"
	"log"

	"github.com/gokapaya/cshelper/errlist"
	"github.com/gokapaya/cshelper/match"
)

const PairBucket = "pairs"

// StoreMatches saves Pair data to the database.
// `pl` can be single or multiple Pair(s)
func (c *Client) StoreMatches(pl ...match.Pair) error {
	var el errlist.ErrList

	tx, err := c.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	b, err := tx.CreateBucketIfNotExists([]byte(PairBucket))
	if err != nil {
		return err
	}

	var total int
	for _, p := range pl {

		value, err := json.Marshal(p)
		if err != nil {
			el = append(el, err)
			continue
		}
		if err := b.Put([]byte(p.Santa.Username), value); err != nil {
			el = append(el, err)
			continue
		}
		total++
	}
	if el.NotEmpty() {
		return el
	}

	log.Printf("Added %v new pairs\n", total)

	return tx.Commit()
}

// GetPair retrieves a single Pair from the database
func (c *Client) GetPair(key string) (*match.Pair, error) {
	tx, err := c.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	b := tx.Bucket([]byte(PairBucket))
	if b == nil {
		return nil, errlist.Err("Bucket not found")
	}

	value := b.Get([]byte(key))
	if value == nil {
		return nil, errlist.Err("Pair not found")
	}

	var p match.Pair
	if err := json.Unmarshal(value, &p); err != nil {
		return nil, err
	}

	return &p, nil
}

// GetMatches retrieves the list of Pairs from the database
func (c *Client) GetMatches() ([]match.Pair, error) {
	var (
		el errlist.ErrList
		pl []match.Pair
	)

	tx, err := c.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	b := tx.Bucket([]byte(PairBucket))
	if b == nil {
		return nil, errlist.Err("Bucket not found")
	}

	cr := b.Cursor()
	for k, v := cr.First(); k != nil; k, v = cr.Next() {
		var p match.Pair
		if err := json.Unmarshal(v, &p); err != nil {
			el = append(el, err)
			continue
		}
		pl = append(pl, p)
	}

	if el.NotEmpty() {
		return pl, el
	}

	return pl, nil
}
