package handlers

import (
	"context"
	"fmt"
)

type IDemo interface {
	DemoByID(ctx context.Context, request interface{}) (interface{}, error)
	Demo(ctx context.Context, request interface{}) (interface{}, error)
}

func NewDemo() IDemo {
	return &DemoValidator{
		next: newDemo(),
	}
}

type demo struct {
}

func newDemo() *demo {
	return &demo{}
}

// Demo return a demo
func (h demo) Demo(ctx context.Context, request interface{}) (interface{}, error) {
	return map[string]string{"demo": "this is demo"}, nil
}

// DemoByID return a demo
func (h demo) DemoByID(ctx context.Context, request interface{}) (interface{}, error) {
	fmt.Println(request)
	return map[string]string{"demo": "this is demo"}, nil
}
