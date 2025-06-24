package main

import (
	"topviewgame/rand"
)

var rnd = rand.New(1 << 23)

func GetRandomInt(num int) int {
	return rnd.Intn(num)
}

// GetDiceRoll returns an integer from 1 to the number
func GetDiceRoll(num int) int {
	return rnd.Intn(num) + 1
}
