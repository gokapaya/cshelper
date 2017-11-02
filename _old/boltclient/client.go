package boltclient

import (
	"time"

	"github.com/boltdb/bolt"
)

// Client represents a client connection
// to a bolt database
type Client struct {
	Path string

	buckets []string
	db      *bolt.DB
}

// NewClient takes a path string and
// returns a new `Client` to that database
func NewClient(path string, bucket ...string) *Client {
	return &Client{
		Path:    path,
		buckets: bucket,
	}
}

// SetBucket registers a bucket (or multiple) with the client
func (c *Client) SetBucket(buckets ...string) {
	for _, b := range buckets {
		c.buckets = append(c.buckets, b)
	}
}

// Open opens a new bolt database and initializes the `buckets`
func (c *Client) Open() error {
	// Open the database
	db, err := bolt.Open(c.Path, 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	c.db = db

	tx, err := c.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// create all the clients buckets
	for _, b := range c.buckets {
		if _, err := tx.CreateBucketIfNotExists([]byte(b)); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Close closes the database connection
func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}
