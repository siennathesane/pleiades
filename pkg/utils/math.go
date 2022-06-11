package utils

import (
	"math/rand"
	"time"
)

func RandomInt(min int, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	p := r.Perm(max - min + 1)
	return p[min]
}
