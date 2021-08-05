package middleware

import (
	"reflect"
	"unsafe"

	"github.com/bcowtech/worker-redis/internal"
)

func isMessageHandler(rv reflect.Value) bool {
	if rv.IsValid() {
		return rv.Type().AssignableTo(typeOfMessageHandler)
	}
	return false
}

func asMessageHandler(rv reflect.Value) internal.MessageHandler {
	if rv.IsValid() {
		if v, ok := rv.Convert(typeOfMessageHandler).Interface().(internal.MessageHandler); ok {
			return v
		}
	}
	return nil
}

func asRedisWorker(rv reflect.Value) *internal.RedisWorker {
	return reflect.NewAt(typeOfHost, unsafe.Pointer(rv.Pointer())).
		Interface().(*internal.RedisWorker)
}
