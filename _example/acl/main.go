package main

import (
	"log"
)

type Permission uint64

const (
	Read Permission = 1 << iota
	Create
	Update
	Delete
)

func main() {
	var ReadAndUpdate Permission = 5
	var hasPermission bool = false

	if ReadAndUpdate == (Read | Create | Update) {
		hasPermission = true
	}
	log.Println("Permission Read and Update: ", ReadAndUpdate)
	log.Println("Has Permission: ", hasPermission)
}
