package redis

import (
	"math/rand"
	"time"
)

func RandomDay() time.Duration {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	return time.Hour * time.Duration(24+r.Int31n(25))
}
