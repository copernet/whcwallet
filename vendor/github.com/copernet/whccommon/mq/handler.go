package mq

import (
	"context"
	"errors"
	"time"

	"github.com/copernet/whccommon/log"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

func BuildFactory(arg interface{}) (*RedisQueueHandler, error) {
	if pool, ok := arg.(*redis.Pool); ok {
		return &RedisQueueHandler{pool: pool}, nil
	}

	return nil, errors.New("BuildFactory Error,find nil")
}

type MessageQueueHandler interface {
	Publish(channel string, msg string, ctx context.Context) (error)
	Subscribe(channel string, consumer ConsumerHandler) error
}

type ConsumerHandler interface {
	OnMessage(msg string, ctx context.Context)
}

type RedisQueueHandler struct {
	pool *redis.Pool
}

func (r *RedisQueueHandler) Publish(channel string, msg string, ctx context.Context) error {
	ri := r.pool.Get()
	defer ri.Close()

	_, err := ri.Do("PUBLISH", channel, msg)
	if err != nil {
		log.WithCtx(ctx).Errorf("publish message:%s to channel:%s error:%s", msg, channel, err.Error())
		return err
	}

	log.WithCtx(ctx).Infof("publish message:%s to channel:%s", msg, channel)
	return err
}

func (r *RedisQueueHandler) Subscribe(channel string, consumer ConsumerHandler) error {
	c := make(chan []byte)
	go func() {
		for balChan := range c {
			consumer.OnMessage(string(balChan), log.NewContext())
		}
	}()

	err := subscribeChannel(context.Background(), channel, c, r.pool)
	if err != nil {
		logrus.Errorf("subscribe chananel for wormhole balance error: %s", err)
		return err
	}

	return nil
}

func subscribeChannel(ctx context.Context, channel string, ret chan<- []byte, pool *redis.Pool) error {
	ri := pool.Get()
	defer ri.Close()

	// A ping is set to the server with this period to test for the health of
	// the connection and server.
	const healthCheckPeriod = time.Minute

	psc := redis.PubSubConn{Conn: ri}

	if err := psc.Subscribe(channel); err != nil {
		return err
	}

	done := make(chan error, 1)

	// Start a goroutine to receive notifications from the server.
	go func() {
		for {
			switch n := psc.Receive().(type) {
			case error:
				done <- n
				return
			case redis.Message:
				ret <- n.Data
			case redis.Subscription:
				if n.Count == 0 {
					// Return from the goroutine when all channels are unsubscribed.
					done <- nil
					return
				}
			}
		}
	}()

	ticker := time.NewTicker(healthCheckPeriod)
	defer ticker.Stop()

	var err error
loop:
	for err == nil {
		select {
		case <-ticker.C:
			// Send ping to test health of connection and server. If
			// corresponding pong is not received, then receive on the
			// connection will timeout and the receive goroutine will exit.
			if err = psc.Ping(""); err != nil {
				break loop
			}
		case <-ctx.Done():
			break loop
		case err := <-done:
			// Return error from the receive goroutine.
			return err
		}
	}

	// Signal the receiving goroutine to exit by unsubscribing from all channels.
	psc.Unsubscribe()

	// Wait for goroutine to complete.
	return <-done

	return nil
}
