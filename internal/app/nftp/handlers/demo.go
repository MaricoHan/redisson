package handlers

import (
	"context"
)

type IDemo interface {
	Demo(ctx context.Context, request interface{}) (interface{}, error)
}

func NewDemo() IDemo {
	return newDemo()
}

type demo struct {
	base
}

func newDemo() *demo {
	return &demo{}
}

// Demo return a demo
func (h demo) Demo(ctx context.Context, request interface{}) (interface{}, error) {
	return map[string]string{"demo": "this is demo"}, nil
}
