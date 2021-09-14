package redis

import (
	redis "github.com/bcowtech/lib-redis-stream"
	"github.com/bcowtech/worker-redis/internal"
)

const (
	AutoIncrement        = redis.AutoIncrement
	LastStreamOffset     = redis.LastStreamOffset
	NextStreamOffset     = redis.NextStreamOffset
	StartingStreamOffset = redis.StartingStreamOffset
)

type (
	Options  = redis.Options
	XMessage = redis.XMessage
	XStream  = redis.XStream

	AdminClient     = redis.AdminClient
	Forwarder       = redis.Forwarder
	ForwarderRunner = redis.ForwarderRunner
	WorkerContext   = redis.ConsumeContext

	MessageHandler = internal.MessageHandler
	Worker         = internal.RedisWorker
)

type (
	RedisErrorHandleProc = internal.RedisErrorHandleProc
)
