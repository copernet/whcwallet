package ws

import (
	"encoding/json"
	"sync"

	"github.com/copernet/whccommon/log"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type ConnMgr struct {
	addrMtx sync.Mutex
	// only accept bech32 encoded address
	connMapping map[string]*Connection

	MsgChan chan Message
}

// addAddrs only accept bech32 form encoded address
func (c *ConnMgr) addAddrs(conn *Connection, addrs []string) {
	c.addrMtx.Lock()
	willRemoveConnections := make([]*Connection, 0)
	// all address point to the same *Connection
	for _, addr := range addrs {
		if connection, ok := c.connMapping[addr]; ok {
			willRemoveConnections = append(willRemoveConnections, connection)

			willRemovedAddrs := make([]string, 0, connection.addrs.Size())
			connection.addrs.Each(func(item interface{}) bool {
				willRemovedAddrs = append(willRemovedAddrs, item.(string))

				return true
			})
			logrus.WithField(log.DefaultTraceLabel, connection.Uid).
				Infof("The account login two devices, hold addresses: %v", willRemovedAddrs)

			logrus.WithField(log.DefaultTraceLabel, conn.Uid).
				Debugf("the new replaced connection for address: %s", addr)
		}

		c.connMapping[addr] = conn
	}
	c.addrMtx.Unlock()

	for _, conn := range willRemoveConnections {
		conn.Close()
	}
}

func (c *ConnMgr) delAddr(addrs []string) {
	c.addrMtx.Lock()
	defer c.addrMtx.Unlock()

	for _, addr := range addrs {
		connection, ok := c.connMapping[addr]
		if ok {
			logrus.WithField(log.DefaultTraceLabel, connection.Uid).
				Infof("remove address from websocket sync: %s", addr)
		}

		delete(c.connMapping, addr)
	}
}

func (c *ConnMgr) GetConnForAddress(addr string) (*Connection, bool) {
	c.addrMtx.Lock()
	defer c.addrMtx.Unlock()

	if conn, ok := c.connMapping[addr]; ok {
		return conn, true
	}

	return nil, false
}

// Notify serves the client for the relative addresses's balance for connection.
// Ignore the unknown address's messages only to log it. The client will be closed
// immediately if writing to client failed.
// The client should be reconnect to the server with all the specified addresses in
// the account.
func (c *ConnMgr) Notify() {
	defer func() {
		logrus.Error("serious problem: balance will not notification. And /transaction/push API will block")
	}()

	for msg := range c.MsgChan {
		connection, ok := c.GetConnForAddress(msg.Address)
		if !ok {
			logrus.Warnf("to notify unknown address balance, address: %s", msg.Address)

			// maybe server mismatch, just continue to serve
			continue
		}

		if msg.Symbol == CoinWormhole {
			logrus.WithField(log.DefaultTraceLabel, connection.Uid).
				Infof("omit wormhole balance for address: %s, balance: %v", msg.Address, msg.Balance)
		} else if msg.Symbol == CoinBch {
			logrus.WithField(log.DefaultTraceLabel, connection.Uid).
				Infof("omit bch balance for address: %s, balance: %v", msg.Address, msg.Balance)
		}

		// should set write message timeout parameter
		content, err := json.Marshal(&msg)
		if err != nil {
			logrus.WithField(log.DefaultTraceLabel, connection.Uid).
				Errorf("marshal Message failed: %s", err)
			continue
		}

		err = connection.Conn.WriteMessage(websocket.TextMessage, content)
		if err != nil {
			logrus.WithField(log.DefaultTraceLabel, connection.Uid).
				Warnf("write balance message to client via websocket failed: %v", err)

			connection.Close()
		}
	}
}

func NewConnMgr() *ConnMgr {
	return &ConnMgr{
		connMapping: make(map[string]*Connection),
		MsgChan:     make(chan Message, 100),
	}
}
