package utils

import (
	"sync"
	"testing"
)

func TestGoID(t *testing.T) {
	w := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		w.Add(1)
		go func() {
			w.Done()
			t.Log(GoID())
		}()
	}
	w.Wait()
}
