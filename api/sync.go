package api

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/bcext/cashutil"
	"github.com/copernet/go-electrum/electrum"
	"github.com/copernet/whccommon/log"
	common "github.com/copernet/whccommon/model"
	"github.com/copernet/whcwallet/config"
	"github.com/copernet/whcwallet/logic/ws"
	"github.com/copernet/whcwallet/model"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

const (
	pingInterval = 5 * time.Second

	// Shutdown if ping to server failed five times continuously.
	// Because the current program rely on electrum server seriously.
	// This system will not supply any most services if connected electrum
	// server work ill.
	retryLimit = 5

	// max retry times fro subscribing bitcoin cash balance update
	retrySubscribe = 3
)

var (
	// prevent from creating two more connection instances
	create sync.Mutex
	node   *electrum.Node
)

func getNode() *electrum.Node {
	create.Lock()
	defer create.Unlock()

	if node != nil && node.Ping() == nil {
		return node
	}

	n := electrum.NewNode()
	if err := n.ConnectTCP(config.GetConf().Electron.Host + ":" + config.GetConf().Electron.Port); err != nil {
		logrus.Errorf("create connection to electrum error: %v", err)
		// unnecessary to star the application if electrum connection failed
		//os.Exit(1)
	}

	node = n
	return node
}

func init() {
	// debug mode. close this option in production environment
	electrum.DebugMode = true

	// keep connected to electrum server
	go keepAlive()

	tipChan, err := getNode().BlockchainHeadersSubscribe()
	if err != nil {
		logrus.Errorf("subscribe blockchain tip update from electrum error: %s", err)
	}

	go func() {
		for tip := range tipChan {
			ctx := log.NewContext()
			log.WithCtx(ctx).Infof("block tip:%s,height:%d", model.UpdateBlockTip, tip.BlockHeight)

			err := model.Publish(model.UpdateBlockTip, "tip updated")
			//logic.DumpHeaders(ctx, tip)

			if err != nil {
				logrus.Error("notify whcengine to update blockchain failed")
			}
		}
	}()

	c := make(chan []byte)

	go func() {
		err := model.Subscribe(context.Background(), model.BalanceUpdated, c)
		if err != nil {
			logrus.Errorf("subscribe channel for wormhole balance error: %s", err)

			// panic to exit program because of initial failed.
			panic(err)
		}
	}()

	go func() {
		for balMsg := range c {
			ctx := log.NewContext()
			log.WithCtx(ctx).Infof("Receive Message:%s from channel:", balMsg)

			var balNotify common.BalanceNotify
			err := json.Unmarshal(balMsg, &balNotify)
			if err != nil {
				log.WithCtx(ctx).Errorf("json unmarshal from redis balance notify failed: %v", err)
				continue
			}

			address, err := cashutil.DecodeAddress(balNotify.Address, config.GetChainParam())
			if err != nil {
				log.WithCtx(ctx).Errorf("mismatch bitcoin address format: %v, address: %s", err, balNotify.Address)
				continue
			}

			addr := address.EncodeAddress(true)
			// prepare for notification
			ret, err := model.GetNotificationItem(addr, balNotify.PropertyID, balNotify.TxID)
			if err != nil {
				if gorm.IsRecordNotFoundError(err) {
					log.WithCtx(ctx).Warnf("notification item not found: address: %s", addr)
				} else {
					log.WithCtx(ctx).Errorf("get notification item failed: %v, %s", err, addr)
				}
			} else {
				err = model.StoreNotification(addr, ret)
				if err != nil {
					log.WithCtx(ctx).Errorf("store notification item failed: %v, address: %s", err, addr)
				}
			}

			connection, ok := connMgr.GetConnForAddress(addr)
			// the client holds this address not online.
			if !ok {
				log.WithCtx(ctx).Warnf("the client holds this address not online, address: %s", addr)
				continue
			}

			bal, err := GetBalanceFromCache([]string{addr}, true)
			if err != nil {
				log.WithCtx(ctx).WithField(log.DefaultTraceLabel, connection.Uid).
					Errorf("get wormhole balance for address: %s error: %v", addr, err)
				continue
			}

			var msg ws.Message
			msg.Address = addr
			msg.Symbol = ws.CoinWormhole
			msg.Balance = bal[addr]

			log.WithCtx(ctx).WithFields(logrus.Fields{
				log.DefaultTraceLabel: connection.Uid,
			}).Debugf("Subscribe: receive wormhole balance update message from redis, address: %s; balance: %v", msg.Address, msg.Balance)
			// notify websocket
			connMgr.MsgChan <- msg
		}
	}()
}

func SyncBalanceForBCH(addr cashutil.Address) {
	base58Addr := addr.EncodeAddress(false)
	bench32Addr := addr.EncodeAddress(true)
	connection, ok := connMgr.GetConnForAddress(bench32Addr)
	// the client holds this address not online.
	if !ok {
		logrus.Infof("the client connection has closed, address: %s", bench32Addr)

		return
	}

	node := getNode()
	balChan, err := node.BlockchainAddressSubscribe(base58Addr)
	if err != nil {
		// retry for subscribing bitcoin cash balance update notification.
		for i := 0; i < retrySubscribe; i++ {
			balChan, err = node.BlockchainAddressSubscribe(base58Addr)
			if err == nil {
				break
			}
		}

		if err != nil {
			logrus.WithFields(logrus.Fields{
				log.DefaultTraceLabel: connection.Uid,
			}).Errorf("track bch balance failed: %v; address: %s", err, bench32Addr)

			return
		}
	}

	logrus.WithFields(logrus.Fields{
		log.DefaultTraceLabel: connection.Uid,
	}).Infof("start track bch balance for address: %s", bench32Addr)

	go func() {
		defer func() {
			connection.Close()

			logrus.WithFields(logrus.Fields{
				log.DefaultTraceLabel: connection.Uid,
			}).Infof("goroutine exit for balance track: %s", bench32Addr)
		}()

		for {
			select {
			case <-connection.Quit:
				logrus.WithFields(logrus.Fields{
					log.DefaultTraceLabel: connection.Uid,
				}).Infof("stop track bch balance for address: %s", bench32Addr)

				return
			case <-balChan:
				bal, err := getNode().BlockchainAddressGetBalance(base58Addr)
				if err != nil {
					// only log this error if encountering some error.
					// It will request again when the coin confirmed.
					// keep the goroutine running.
					logrus.WithFields(logrus.Fields{
						log.DefaultTraceLabel: connection.Uid,
					}).Errorf("get BCH balance from electrum server error: %v; address: %s", err, bench32Addr)

					continue
				}

				// assemble websocket message
				var msg ws.Message
				msg.Address = bench32Addr
				msg.Balance = ws.BCHBalance{
					Confirmed:   bal.Confirmed.ToBCH(),
					Unconfirmed: bal.Unconfirmed.ToBCH(),
				}
				msg.Symbol = ws.CoinBch

				logrus.WithFields(logrus.Fields{
					log.DefaultTraceLabel: connection.Uid,
				}).Debugf("Subscribe: receive bitcoin-cash balance message from electrum server, address: %s; balance: %v", bench32Addr, msg.Balance)

				// update websocket
				connMgr.MsgChan <- msg
			}
		}
	}()
}

func keepAlive() {
	var reties int

	for {
		if err := getNode().Ping(); err != nil {
			reties++
			if reties >= retryLimit {
				logrus.Error("Retry to connect to electrum server failed too many times")

				// Should find a HA resolution before stop exit program directly
				os.Exit(1)
			}

			logrus.Errorf("Ping to electrum server error: %v", err)
		} else if reties != 0 {
			reties = 0
		}

		time.Sleep(pingInterval)
	}
}
