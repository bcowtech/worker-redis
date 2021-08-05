package redis

import (
	"github.com/bcowtech/host"
	"github.com/bcowtech/worker-redis/internal/middleware"
)

func UseErrorHandler(handler RedisErrorHandleProc) host.Middleware {
	if handler == nil {
		panic("argument 'handler' cannot be nil")
	}

	return &middleware.ErrorHandlerMiddleware{
		Handler: handler,
	}
}

func UseStreamGateway(streamGateway interface{}) host.Middleware {
	if streamGateway == nil {
		panic("argument 'topicGateway' cannot be nil")
	}

	return &middleware.StreamGatewayMiddleware{
		StreamGateway: streamGateway,
	}
}
