package middleware

import (
	"reflect"

	"github.com/bcowtech/worker-redis/internal"
)

const (
	UNHANDLED_MESSAGE_HANDLER_TOPIC_SYMBOL = "?"
)

var (
	typeOfHost           = reflect.TypeOf(internal.RedisWorker{})
	typeOfMessageHandler = reflect.TypeOf((*internal.MessageHandler)(nil)).Elem()

	TAG_STREAM = "stream"
	TAG_OFFSET = "offset"
)

type (
	ConfigureUnhandledMessageHandleProc func(handler internal.MessageHandler)
	ConfigureStream                     func(stream internal.StreamOffset)
)
