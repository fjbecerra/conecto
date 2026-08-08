package shared

import (
	"conecto/core/retry"
	"conecto/shared/config"
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

func NewRetry() Retry{
	return Retry{
		Random: &RandomImpl{},
	} 

}

type Retry struct{
	Random retry.Random
}

func (r * Retry) CreateRetryExecutor(config *config.Retry) retry.Executor {
	var retryPolicy retry.Policy
	if config != nil {
		retryPolicy = retry.Policy{
			MaxRetries:     config.MaxRetries,
			InitialBackoff: time.Duration(config.BackoffMS),
			MaxBackoff:     time.Duration(config.MaxBackoff),
			Jitter:         true,
		}
	}else{
		retryPolicy = retry.Policy{
			MaxRetries:     3,
			InitialBackoff: time.Duration(500),
			MaxBackoff:     time.Duration(500),
			Jitter:         true,
		}
	}
	return retry.Executor{
		Policy: retryPolicy,
		Random: r.Random,
	}
}