package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/metric"
	"gitlab.bianjie.ai/avata/open-api/utils"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
)

type grpcInterceptorMiddleware struct {
}

var (
	grpcInterceptorMiddlewares    *grpcInterceptorMiddleware
	grpcInterceptorMiddlewareOnce sync.Once
)

// NewGrpcInterceptorMiddleware 单列模式
func NewGrpcInterceptorMiddleware() *grpcInterceptorMiddleware {
	grpcInterceptorMiddlewareOnce.Do(func() {
		grpcInterceptorMiddlewares = &grpcInterceptorMiddleware{}
	})
	return grpcInterceptorMiddlewares
}

// Interceptor grpc 拦截器
// ctx context.Context是请求上下文
// cc *ClientConn是调用RPC的客户端连接
// method string是请求的方法名
// opts ...CallOption包含了所有适用呼叫选项，包括来自于客户端连接的默认选项和所有的呼叫。
func (g *grpcInterceptorMiddleware) Interceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		authData, err := g.authData(ctx)
		if err != nil {
			log.WithError(err).Errorln("interceptor")
			return err
		}
		status := "200"
		err = invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			status = g.handleErrorCodeToString(err)
		}
		metric.NewPrometheus().ApiServiceRequests.With([]string{
			"name", fmt.Sprintf("%s-%s", authData.Code, authData.Module),
			"method", method,
			"status", status,
		}...).Add(1)
		return err
	}
}

// authData 解析Context中的数据
func (g *grpcInterceptorMiddleware) authData(ctx context.Context) (vo.AuthData, error) {
	authDataString := ctx.Value("X-Auth-Data")
	authDataSlice, ok := authDataString.([]string)
	if !ok {
		return vo.AuthData{}, fmt.Errorf("the key doesnt exist")
	}
	var authData vo.AuthData
	err := json.Unmarshal([]byte(authDataSlice[0]), &authData)
	if err != nil {
		return vo.AuthData{}, fmt.Errorf("json un marshal error")
	}
	return authData, nil
}

// handleErrorCodeToString 解析异常中的Code
func (g *grpcInterceptorMiddleware) handleErrorCodeToString(err error) string {
	respErr := errors2.Convert(err)
	code := "500"
	if utils.IsNumeric(respErr.Code().String()[5:8]) {
		code = respErr.Code().String()[5:8]
	}
	return code
}
