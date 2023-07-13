package random

import (
	"math/rand"
	"time"
)

func NewRandomString(length int) string {
	var gen string
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for len(gen) < length {
		num := rnd.Int31n(26)
		lowercase := rnd.Int31n(2)
		gen += string(65 + num + lowercase*32)
	}

	return gen
}
