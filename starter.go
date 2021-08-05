package redis

import (
	"github.com/bcowtech/host"
	"github.com/bcowtech/worker-redis/internal"
)

func Startup(app interface{}, middlewares ...host.Middleware) *host.Starter {
	starter := host.Startup(app, middlewares...)
	// register HostProvider
	host.RegisterHostService(starter, internal.RedisWorkerServiceInstance)

	return starter
}
