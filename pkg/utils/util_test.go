package utils

import (
	"fmt"
	"testing"
)

func TestGoID(t *testing.T) {
	for i := 0; i < 100; i++ {
		go func() {
			fmt.Println(GoID())
		}()
	}
}
