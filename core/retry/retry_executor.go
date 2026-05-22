package retry

import (
	"context"
	"time"
)


type Executor struct {
    Policy Policy
    Random   Random
}

func (e *Executor) Do(ctx context.Context, fn func() error) error {

    var err error
    backoff := e.Policy.InitialBackoff

    for i := 0; i <= e.Policy.MaxRetries; i++ {

        err = fn()
        if err == nil {
            return nil
        }

        if i == e.Policy.MaxRetries {
            return err
        }

        delay := backoff

        if e.Policy.Jitter {
            delay = jitter(delay, e.Random.Float64)
        }

        if delay > e.Policy.MaxBackoff {
            delay = e.Policy.MaxBackoff
        }

        select {
        case <-time.After(delay):
        case <-ctx.Done():
            return ctx.Err()
        }

        backoff = backoff * 2
    }

    return err
}

func jitter(d time.Duration, randFn func() float64) time.Duration {
    j := float64(d) * (0.5 + randFn())
    return time.Duration(j)
}

