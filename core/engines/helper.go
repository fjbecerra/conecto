package engines

import (
	"math/rand/v2"
	"time"
)

func BackoffWithJitter(base time.Duration, attempt int, max time.Duration, r *rand.Rand) time.Duration {
	d := base * time.Duration(1<<attempt)

	if d > max {
		d = max
	}
	if d <= 0 {
		return base
	}

	return time.Duration(r.Int64N(int64(d)))
}