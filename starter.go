package redis

import (
	"github.com/bcowtech/host"
	"github.com/bcowtech/worker-redis/internal"
)

func Startup(app interface{}) *host.Starter {
	var (
		starter = host.Startup(app)
	)

	host.RegisterHostService(starter, internal.RedisWorkerServiceInstance)

	return starter
}
