package handlers

import (
	"context"
)

type DemoRequest struct {
	ID uint `json:"id"`
}

type DemoValidator struct {
	next IDemo
}

// Demo return a demo
func (h DemoValidator) Demo(ctx context.Context, _ interface{}) (interface{}, error) {
	//  验证参数
	return h.next.Demo(ctx, nil)
}

// DemoByID return a demo
func (h DemoValidator) DemoByID(ctx context.Context, _ interface{}) (interface{}, error) {
	//  验证参数
	request := DemoRequest{
		ID: 1,
	}
	return h.next.DemoByID(ctx, request)
}
