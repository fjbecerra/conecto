package main

import (
	"fmt"
	"testing"
)


func TestGoRutines(t *testing.T) {
	ch := make(chan int)

    worker(ch)

    ch <- 1
    ch <- 2
    ch <- 3

    close(ch)
}

func worker(ch chan int) {
    for n := range ch {
        fmt.Println("Processing:", n)
    }
}