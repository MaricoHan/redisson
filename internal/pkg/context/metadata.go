package context

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// SetMetadataToContext 设置元数据到Context
func SetMetadataToContext(ctx context.Context, val map[string]string) context.Context {
	md := metadata.New(val)
	return metadata.NewIncomingContext(ctx, md)
}
