package util

import (
	"fmt"
	"testing"
)

func TestGoID(t *testing.T) {
	for i := 0; i < 100; i++ {
		fmt.Println(GoID())
	}
}
