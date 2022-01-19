package handlers

import "context"

//type IBase interface {
//	AppID(ctx context.Context) uint64
//}

type base struct {
}

func (h base) AppID(ctx context.Context) uint64 {
	appID := ctx.Value("X-App-ID")
	return appID.(uint64)
}
