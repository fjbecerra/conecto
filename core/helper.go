package core

import (
	"context"
	"encoding/json"
	"os"
)

func StreamMap[T any, R any](
    ctx context.Context,
    in <-chan T,
    fn func(T) (R, error),
) (<-chan R, <-chan error) {

    out := make(chan R)
    errs := make(chan error, 1)

    go func() {
        defer close(out)
        defer close(errs)

        for item := range in {
            select {
				case <-ctx.Done():
					errs <- ctx.Err()
					return
				default:
            }

            result, err := fn(item)
            if err != nil {
                errs <- err
                return
            }

            out <- result
        }
    }()

    return out, errs
}

func LoadConfig[T any](path string) T{
    data, err := os.ReadFile(path)
    if err != nil {
        panic(err)
    }
    var config T

    err = json.Unmarshal(data, &config)
    if err != nil {
        panic(err)
    }
    return config
}