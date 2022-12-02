package main

import (
	"math/rand"
	"time"
)

func getRandomSpaceQuote() string {
	rand.Seed(time.Now().Unix())
	quote := rand.Intn(2)
	switch quote {
	case 0:
		return "“The stars don't look bigger, but they do look brighter.” ― Sally Ride"
	case 1:
		return "“I see Earth! It is so beautiful.” ― Yuri Gagarin"
	default:
		return ""
	}
}
