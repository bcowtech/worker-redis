package standardtest

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/bcowtech/config"
	redis "github.com/bcowtech/worker-redis"
)

func TestStarter(t *testing.T) {
	err := setupTestStarter()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := teardownTestStarter()
		if err != nil {
			t.Fatal(err)
		}
	}()

	app := App{}
	starter := redis.Startup(&app).
		Middlewares(
			redis.UseStreamGateway(&StreamGateway{}),
			redis.UseErrorHandler(func(err error) (disposed bool) {
				t.Logf("catch err: %v", err)
				return false
			}),
		).
		ConfigureConfiguration(func(service *config.ConfigurationService) {
			service.
				LoadEnvironmentVariables("").
				LoadYamlFile("config.yaml").
				LoadCommandArguments()
		})

	runCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := starter.Start(runCtx); err != nil {
		t.Error(err)
	}

	select {
	case <-runCtx.Done():
		if err := starter.Stop(context.Background()); err != nil {
			t.Error(err)
		}
	}

	// assert app.Config
	{
		conf := app.Config
		if conf.RedisAddress == "" {
			t.Errorf("assert 'Config.RedisAddress':: should not be an empty string")
		}
		var expectedRedisConsumerGroup string = "gotestGroup"
		if conf.RedisConsumerGroup != expectedRedisConsumerGroup {
			t.Errorf("assert 'Config.RedisConsumerGroup':: expected '%v', got '%v'", expectedRedisConsumerGroup, conf.RedisConsumerGroup)
		}
		var expectedRedisConsumerName string = "gotestConsumer"
		if conf.RedisConsumerName != expectedRedisConsumerName {
			t.Errorf("assert 'Config.RedisConsumerName':: expected '%v', got '%v'", expectedRedisConsumerName, conf.RedisConsumerName)
		}
		var expectedRedisMaxInFlight int64 = 8
		if conf.RedisMaxInFlight != expectedRedisMaxInFlight {
			t.Errorf("assert 'Config.RedisMaxInFlight':: expected '%v', got '%v'", expectedRedisMaxInFlight, conf.RedisMaxInFlight)
		}
		var expectedRedisMaxPollingTimeout time.Duration = 10 * time.Millisecond
		if conf.RedisMaxPollingTimeout != expectedRedisMaxPollingTimeout {
			t.Errorf("assert 'Config.RedisMaxPollingTimeout':: expected '%v', got '%v'", expectedRedisMaxPollingTimeout, conf.RedisMaxPollingTimeout)
		}
		var expectedRedisAutoClaimMinIdleTime time.Duration = 30 * time.Second
		if conf.RedisAutoClaimMinIdleTime != expectedRedisAutoClaimMinIdleTime {
			t.Errorf("assert 'Config.RedisAutoClaimMinIdleTime':: expected '%v', got '%v'", expectedRedisAutoClaimMinIdleTime, conf.RedisAutoClaimMinIdleTime)
		}
	}
}

func setupTestStarter() error {
	opt := &redis.Options{
		Addr: os.Getenv("REDIS_SERVER"),
		DB:   0,
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
		_, err = admin.CreateConsumerGroupAndStream("gotestStream", "gotestGroup", redis.LastStreamOffset)
		if err != nil {
			return err
		}

		p, err := redis.NewForwarder(opt)
		if err != nil {
			return err
		}
		defer p.Close()

		_, err = p.Write("gotestStream", redis.AutoIncrement, map[string]interface{}{
			"name": "luffy",
			"age":  19,
		})
		if err != nil {
			return err
		}
		_, err = p.Write("gotestStream", redis.AutoIncrement, map[string]interface{}{
			"name": "nami",
			"age":  21,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func teardownTestStarter() error {
	admin, err := redis.NewAdminClient(&redis.Options{
		Addr: os.Getenv("REDIS_SERVER"),
		DB:   0,
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
