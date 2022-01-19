package internal

import (
	"context"
	"log"
	"os"
	"sync/atomic"
	"testing"
	"time"

	redis "github.com/bcowtech/lib-redis-stream"
)

func TestRedisWorker(t *testing.T) {
	var err error
	err = setupTestRedisWorker()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := teardownTestRedisWorker()
		if err != nil {
			t.Fatal(err)
		}
	}()

	opt := redis.UniversalOptions{
		Addrs: []string{os.Getenv("REDIS_SERVER")},
		DB:    0,
	}

	worker := &RedisWorker{
		ConsumerGroup:     "gotestGroup",
		ConsumerName:      "gotestConsumer",
		RedisOption:       &opt,
		MaxInFlight:       8,
		MaxPollingTimeout: 10 * time.Millisecond,
		ClaimMinIdleTime:  30 * time.Millisecond,
	}

	handler := new(mockMessageHandler)
	worker.preInit()
	{
		worker.dispatcher.streams["gotestStream"] = StreamOffset{
			Stream: "gotestStream",
			Offset: redis.StreamNeverDeliveredOffset,
		}
		worker.dispatcher.router.Add("gotestStream", handler)
	}
	worker.init()

	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	worker.Start(ctx)

	select {
	case <-ctx.Done():
		var expectedMsgCnt int32 = 2
		if handler.msgCnt != expectedMsgCnt {
			t.Errorf("expect %d messages, but got %d messages", expectedMsgCnt, handler.msgCnt)
		}
		worker.Stop(context.Background())
		break
	}
}

type mockMessageHandler struct {
	msgCnt int32
}

func (h *mockMessageHandler) ProcessMessage(ctx *ConsumeContext, stream string, message *XMessage) {
	log.Printf("Message on %s: %v\n", stream, message)
	ctx.Ack(stream, message.ID)

	atomic.AddInt32(&h.msgCnt, 1)
}

func setupTestRedisWorker() error {
	opt := &redis.UniversalOptions{
		Addrs: []string{os.Getenv("REDIS_SERVER")},
		DB:    0,
	}

	admin, err := redis.NewAdminClient(opt)
	if err != nil {
		return err
	}
	defer admin.Close()

	{
		/*
			DEL gotestStream
		*/
		_, err = admin.Handle().Del("gotestStream").Result()
		if err != nil {
			return err
		}

		/*
			XGROUP CREATE gotestStream gotestGroup $ MKSTREAM

			XADD gotestStream * name luffy age 19
			XADD gotestStream * name nami age 21
		*/
		_, err = admin.CreateConsumerGroupAndStream("gotestStream", "gotestGroup", redis.StreamLastDeliveredID)
		if err != nil {
			return err
		}

		p, err := redis.NewProducer(opt)
		if err != nil {
			return err
		}
		defer p.Close()

		_, err = p.Write("gotestStream", redis.StreamAsteriskID, map[string]interface{}{
			"name": "luffy",
			"age":  19,
		})
		if err != nil {
			return err
		}
		_, err = p.Write("gotestStream", redis.StreamAsteriskID, map[string]interface{}{
			"name": "nami",
			"age":  21,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func teardownTestRedisWorker() error {
	admin, err := redis.NewAdminClient(&redis.UniversalOptions{
		Addrs: []string{os.Getenv("REDIS_SERVER")},
		DB:    0,
	})
	if err != nil {
		return err
	}
	defer admin.Close()

	{
		/*
			XGROUP DESTROY gotestStream gotestGroup
		*/
		_, err = admin.DeleteConsumerGroup("gotestStream", "gotestGroup")
		if err != nil {
			return err
		}

		/*
			DEL gotestStream
		*/
		_, err = admin.Handle().Del("gotestStream").Result()
		if err != nil {
			return err
		}
	}
	return nil
}
