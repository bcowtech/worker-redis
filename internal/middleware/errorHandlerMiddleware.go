package middleware

import (
	"github.com/bcowtech/host"
	"github.com/bcowtech/worker-redis/internal"
)

var _ host.Middleware = new(ErrorHandlerMiddleware)

type ErrorHandlerMiddleware struct {
	Handler internal.RedisErrorHandleProc
}

func (m *ErrorHandlerMiddleware) Init(appCtx *host.AppContext) {
	var (
		kafkaworker = asRedisWorker(appCtx.Host())
		preparer    = internal.NewRedisWorkerPreparer(kafkaworker)
	)

	preparer.RegisterErrorHandler(m.Handler)
}
