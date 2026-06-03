package factories

import (
	"math/rand/v2"
	"time"
)

type RandomImpl struct {}

func (random * RandomImpl) Float64() float64 {
	seed := time.Now().UnixNano()

	r := rand.New(rand.NewPCG(
		uint64(seed),
		uint64(seed>>1),
	))
	return r.Float64()
	
}

