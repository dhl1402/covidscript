package playground

import (
	"context"

	"github.com/dhl1402/covidscript/cmd/api/svc"
	"github.com/go-kit/kit/endpoint"
)

func makeInterpreteEndpoint(s Service, middlewares ...endpoint.Middleware) endpoint.Endpoint {
	e := func(ctx context.Context, req interface{}) (res interface{}, err error) {
		return s.Intepret(ctx, req.(InterpretRequest))
	}
	return svc.NewEndpointWithMiddlewares(e, middlewares...)
}
