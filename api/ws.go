package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/bcext/cashutil"
	"github.com/copernet/whccommon/log"
	"github.com/copernet/whcwallet/config"
	"github.com/copernet/whcwallet/logic/ws"
	"github.com/copernet/whcwallet/util"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

var connMgr = ws.NewConnMgr()

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type action string

const (
	add    action = "add"
	remove action = "remove"
)

// type only select from add/remove
type OperateAddrsMsg struct {
	Type      action   `json:"type"`
	Addresses []string `json:"addresses"`
}

func NotifyBalanceUpdated(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.WithCtx(ctx).Errorf("upgrade to websocket failed: %v", err)
		return
	}

	// generate uuid fro websocket connection
	uniqueID, _ := uuid.NewV4()
	uid := uniqueID.String()
	ctx.Set(log.DefaultTraceLabel, uid)
	log.WithCtx(ctx).Debug("websocket connect incoming")

	conn.SetCloseHandler(func(code int, text string) error {
		message := websocket.FormatCloseMessage(code, "")
		conn.WriteControl(websocket.CloseMessage, message,
			time.Now().Add(time.Second))

		conn.Close()
		return nil
	})

	c := ws.NewConnection(connMgr, conn, uid)
	go func() {
		defer c.Close()

		for {
			select {
			case <-c.Quit:
				logrus.WithFields(logrus.Fields{
					log.DefaultTraceLabel: c.Uid,
				}).Info("stop reading message, the websocket connection closed")

				return
			default:
			}

			_, mes, err := conn.ReadMessage()
			if err != nil {
				log.WithCtx(ctx).Warnf("read from client message failed: %v", err)
				return
			}

			log.WithCtx(ctx).Debugf("receive message from websocket: %s", string(mes))

			var msg OperateAddrsMsg
			err = json.Unmarshal(mes, &msg)
			if err != nil {
				log.WithCtx(ctx).Errorf("read from client message json unmarshal error: %v", err)
				return
			}

			// limit address encode format
			addrs, err := util.ConvToCashAddr(msg.Addresses, config.GetChainParam())
			if err != nil {
				log.WithCtx(ctx).Errorf("address list encode error: %v", err)
				return
			}

			// limit request added address list length
			if len(addrs) > maxRequestAddressList || len(addrs) == 0 {
				log.WithCtx(ctx).Warnf("address list length illegal, total: %d", len(addrs))
				return
			}

			// add addresses
			if msg.Type == add {
				// Must add websocket connection to cache firstly.
				c.AddAddrs(addrs)

				OnAddressAdd(addrs)
				log.WithCtx(ctx).Debug("add address via websocket: ", msg.Addresses)
			}

			// remove addresses
			if msg.Type == remove {
				c.DelAddr(addrs)
				log.WithCtx(ctx).Debug("delete addresses via websocket: ", msg.Addresses)
			}
			// ignore if the Type is other string
		}
	}()

	// ensure the websocket connection is available.
	go func() {
		timer := time.NewTicker(10 * time.Second)
		defer timer.Stop()
		defer c.Close()

		n := 1
		for {
			select {
			case <-c.Quit:
				logrus.WithFields(logrus.Fields{
					log.DefaultTraceLabel: c.Uid,
				}).Info("stop ping, the websocket connection closed")

				return
			case <-timer.C:
				err := conn.WriteMessage(websocket.TextMessage, []byte(`{"ping":`+strconv.Itoa(n)+`}`))
				if err != nil {
					log.WithCtx(ctx).Warnf("ping to client via websocket failed: %v", err)
					return
				}

				n++
			}
		}
	}()
}

func OnAddressAdd(addrs []string) {
	// start track the bch balance for these addresses
	for _, addr := range addrs {
		// ignore the unexpected error, because upstream logic ensure that
		// the addresses is encoded correctly.
		address, _ := cashutil.DecodeAddress(addr, config.GetChainParam())

		// the electrum only recognises base32 format addresses.
		// SyncBalanceForBCH handles the situation the subscribe balance update fail.
		SyncBalanceForBCH(address)
	}
}

func init() {
	go connMgr.Notify()
}
