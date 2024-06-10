package main

import (
	"testing"
	"time"
)

func TestOrChannel(t *testing.T) {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()

	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(2*time.Second),
	)

	const realOut = 1

	timeAfter := int(time.Since(start).Seconds())
	if timeAfter != realOut {
		t.Errorf("%v!=%v\nShould: %d\nGot: %d\n", realOut, timeAfter, realOut, timeAfter)
	}
}
