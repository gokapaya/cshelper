package boltclient

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/gokapaya/cshelper/csv"
	"github.com/gokapaya/cshelper/errlist"
)

const UserBucket = "giftees"

// StoreUsers saves User data to the database.
// `ul` can be single or multiple User(s)
func (c *Client) StoreUsers(ul ...csv.User) error {
	var el errlist.ErrList

	tx, err := c.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	b, err := tx.CreateBucketIfNotExists([]byte(UserBucket))
	if err != nil {
		return err
	}

	var total int
	for _, u := range ul {

		if v := b.Get([]byte(u.Username)); v != nil {
			continue
		}

		value, err := json.Marshal(u)
		if err != nil {
			el = append(el, err)
			continue
		}
		if err := b.Put([]byte(u.Username), value); err != nil {
			el = append(el, err)
			continue
		}
		total++
	}
	if el.NotEmpty() {
		return el
	}

	log.Printf("Added %v new users\n", total)

	return tx.Commit()
}

// GetUser retrieves a single User from the database
func (c *Client) GetUser(key string) (*csv.User, error) {
	tx, err := c.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	b := tx.Bucket([]byte(UserBucket))
	if b == nil {
		return nil, errlist.Err("Bucket not found")
	}

	value := b.Get([]byte(strings.ToLower(key)))
	if value == nil {
		return nil, errlist.Err("User not found")
	}

	var u csv.User
	if err := json.Unmarshal(value, &u); err != nil {
		return nil, err
	}

	return &u, nil
}

// GetUserList retrieves the list of Users from the database
func (c *Client) GetUserList() ([]csv.User, error) {
	var (
		el errlist.ErrList
		ul []csv.User
	)

	tx, err := c.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	b := tx.Bucket([]byte(UserBucket))
	if b == nil {
		return nil, errlist.Err("Bucket not found")
	}

	cr := b.Cursor()
	for k, v := cr.First(); k != nil; k, v = cr.Next() {
		var u csv.User
		if err := json.Unmarshal(v, &u); err != nil {
			el = append(el, err)
			continue
		}
		ul = append(ul, u)
	}

	if el.NotEmpty() {
		return ul, el
	}

	return ul, nil
}
