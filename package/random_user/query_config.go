package main

import (
	"net/url"
	"strconv"
)

const (
	MAX_RESULTS = 5000
	MIN_RESULTS = 1
)

const (
	KEY_RESULTS  = "results"
	KEY_GENDER   = "gender"
	KEY_PASSWORD = "password"
	KEY_SEED     = "seed"
)

type QueryConfig struct {
	MaxResults int
	Gender     string
	Password   string
	Seed       string
}

func NewQueryConfig() *QueryConfig {
	return &QueryConfig{MIN_RESULTS, "", "", ""}
}

func (c *QueryConfig) encode(u *url.URL) {
	q := u.Query()

	if c.MaxResults < 0 {
		c.MaxResults = MIN_RESULTS
	} else if c.MaxResults > MAX_RESULTS {
		c.MaxResults = MAX_RESULTS
	}

	q.Set(KEY_RESULTS, strconv.Itoa(c.MaxResults))

	q.Set(KEY_GENDER, c.Gender)

	q.Set(KEY_PASSWORD, c.Password)

	q.Set(KEY_SEED, c.Seed)

	u.RawQuery = q.Encode()
}
