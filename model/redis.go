package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/bcext/cashutil"
	"github.com/copernet/whccommon/model"
	"github.com/copernet/whcwallet/config"
	"github.com/copernet/whcwallet/model/view"
	"github.com/gomodule/redigo/redis"
	"golang.org/x/net/context"
)

const (
	BalancePrefix = "balance:"

	// WARNING: the maintained address set only supports bash32 encoded address
	AddressSetKey = "address:set"

	// channels
	// update block tip
	UpdateBlockTip = "block:tip"
	MempoolTxTip = "transaction:tip"
	FeeRateKey     = "transaction:feeRate"
	// notify balance for wormhole
	BalanceUpdated = "balance:wormhole:updated"

	// Redis will persist the balance information for 20 minute unless the program
	// refresh it or timeout.
	BalanceExpire = 1200
)

var pool *redis.Pool

func ConnRedis() error {
	conf := config.GetConf()
	p, err := model.ConnectRedis(conf.Redis)
	if err != nil {
		return err
	}

	pool = p.Pool
	return nil
}

func GetBalanceForAddress(addresses []string) (map[string][]model.BalanceForAddress, error) {
	ri := pool.Get()
	defer ri.Close()

	// the balMap is first return result, and initialize its value via make.
	// so the caller should not use `ret == nil` to justify whether the result is empty
	// or not.
	balMap := make(map[string][]model.BalanceForAddress)
	for _, addr := range addresses {
		_, err := cashutil.DecodeAddress(addr, config.GetChainParam())
		if err != nil {
			// skip to the next address.
			continue
		}

		res, err := redis.String(ri.Do("GET", BalancePrefix+addr))

		if err != nil {
			// continue for the next address if null data in redis database
			if strings.Contains(err.Error(), "nil returned") {
				continue
			}

			return nil, err
		}

		var bal []model.BalanceForAddress
		err = json.Unmarshal(bytes.NewBufferString(res).Bytes(), &bal)
		if err != nil {
			return nil, err
		}

		balMap[addr] = bal
	}

	return balMap, nil
}

func StoreBalanceForAddress(bals map[string][]model.BalanceForAddress) error {
	if bals == nil {
		return errors.New("empty balance result")
	}

	ri := pool.Get()
	defer ri.Close()

	for addr, bal := range bals {
		b, err := json.Marshal(bal)
		if err != nil {
			return err
		}

		// use pipline style to optimise redis operation
		ri.Send("SET", BalancePrefix+addr, string(b), "EX", BalanceExpire)
	}

	err := ri.Flush()
	if err != nil {
		return err
	}

	return err
}

func Publish(channel string, msg string) error {
	ri := pool.Get()
	defer ri.Close()

	_, err := ri.Do("PUBLISH", channel, msg)

	return err
}

func Subscribe(ctx context.Context, channel string, ret chan<- []byte) error {
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
}

func GetFeeRate() (*view.FeeRate, error) {
	ri := pool.Get()
	defer ri.Close()

	res, err := redis.String(ri.Do("GET", FeeRateKey))
	if err != nil {
		return nil, err
	}

	var feeRate view.FeeRate
	err = json.Unmarshal(bytes.NewBufferString(res).Bytes(), &feeRate)
	if err != nil {
		return nil, err
	}

	return &feeRate, nil
}

func StoreFeeRate(fee *view.FeeRate) error {
	ri := pool.Get()
	defer ri.Close()

	b, err := json.Marshal(fee)
	if err != nil {
		return err
	}

	err = ri.Send("SET", FeeRateKey, string(b), "EX", BalanceExpire)
	if err != nil {
		return err
	}

	return nil
}
