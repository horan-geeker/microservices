package util

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandomString returns a random string of length n.
func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// RandomN returns a random number with n digits.
func RandomN(n int) int {
	if n <= 0 {
		return 0
	}
	// min = 10^(n-1), max = 10^n - 1
	var min, max int = 1, 1
	for i := 0; i < n-1; i++ {
		min *= 10
	}
	for i := 0; i < n; i++ {
		max *= 10
	}
	max = max - 1
	return seededRand.Intn(max-min+1) + min
}
