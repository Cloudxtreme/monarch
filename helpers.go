package main

import (
	"math/rand"
	"time"
)

// spit out a random number between to limits
func random(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}
