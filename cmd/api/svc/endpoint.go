package svc

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"go.uber.org/zap"
)

func NewEndpointWithMiddlewares(e endpoint.Endpoint, middlewares ...endpoint.Middleware) endpoint.Endpoint {
	for _, m := range middlewares {
		e = m(e)
	}
	return e
}

func NewLoggerMiddleware(logger *zap.SugaredLogger, name string) endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				logFunc := logger.Infow
				var errStr string
				if err != nil {
					logFunc = logger.Errorw
					errStr = err.Error()
				}
				logFunc(
					name,
					"request", request,
					"time", time.Since(begin),
					"error", errStr,
				)
			}(time.Now())
			response, err = e(ctx, request)
			return
		}
	}
}
