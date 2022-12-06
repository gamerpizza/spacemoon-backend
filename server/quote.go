package main

import (
	"math/rand"
	"time"
)

func getRandomSpaceQuote() string {
	rand.Seed(time.Now().Unix())
	quote := rand.Intn(3)
	switch quote {
	case 0:
		return "“The stars don't look bigger, but they do look brighter.” ― Sally Ride"
	case 1:
		return "“I see Earth! It is so beautiful.” ― Yuri Gagarin"
	case 2:
		return "“Looking up, I see the immensity of the cosmos; bowing my head, I look at the multitude of the world. " +
			"The gaze flies, the heart expands, the joy of the senses can reach its peak, and indeed, " +
			"this is true happiness.“ ― Samantha Cristoforetti (Translation of a Chinese Poem by Wang Xizhi)"
	default:
		return ""
	}
}
