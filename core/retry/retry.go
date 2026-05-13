package retry

import "time"

type Policy struct {
    MaxRetries int
    InitialBackoff time.Duration
    MaxBackoff     time.Duration
    Jitter bool
}