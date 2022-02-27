package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const RANDOM_USER_API = "https://randomuser.me/api/"

func Generate(c *QueryConfig) ([]User, error) {
	u, err := url.Parse(RANDOM_USER_API)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	c.encode(u)

	response, err := http.Get(RANDOM_USER_API + "?" + u.Query().Encode())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer response.Body.Close()

	rawBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var jsonRandomUser JsonRandomUser
	err = json.Unmarshal(rawBytes, &jsonRandomUser)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return jsonRandomUser.Results, nil
}
