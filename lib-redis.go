package redis

import redis "github.com/bcowtech/lib-redis-stream"

func NewAdminClient(opt *Options) (*AdminClient, error) {
	return redis.NewAdminClient(opt)
}

func NewForwarder(opt *Options) (*Forwarder, error) {
	return redis.NewForwarder(opt)
}
