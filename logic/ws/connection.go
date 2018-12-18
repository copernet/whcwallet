package ws

import (
	"sync"

	"github.com/copernet/whccommon/log"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"gopkg.in/fatih/set.v0"
)

type Connection struct {
	// hold all connections information
	connMgr *ConnMgr
	Conn    *websocket.Conn

	// hold address(string) and only accept bech32 encoded address
	addrs set.Interface

	// for websocket trace log
	Uid string

	QuitMutex sync.Mutex
	Quit      chan struct{}
}

// AddAddrs only accept bech32 encoded address
func (conn *Connection) AddAddrs(addrs []string) {
	for i := 0; i < len(addrs); i++ {
		conn.addrs.Add(addrs[i])
	}

	// maintain the connection manager
	conn.connMgr.addAddrs(conn, addrs)
}

// AddAddrs only accept bech32 encoded address
func (conn *Connection) DelAddr(addrs []string) {
	// convert []string to []interface
	addresses := make([]interface{}, len(addrs))
	for i := 0; i < len(addrs); i++ {
		addresses[i] = addrs[i]
	}
	conn.addrs.Remove(addresses...)

	conn.connMgr.delAddr(addrs)
}

// Close should be only called by ConnMgr in order to
// maintain the ConnMgr state and data.
func (conn *Connection) Close() {
	conn.QuitMutex.Lock()
	defer conn.QuitMutex.Unlock()

	select {
	case <-conn.Quit:
		return
	default:
		logrus.WithFields(logrus.Fields{
			log.DefaultTraceLabel: conn.Uid,
		}).Info("close websocket quit channel")

		close(conn.Quit)
	}

	// reuse the variable(memory)
	var addresses []string
	conn.addrs.Each(func(item interface{}) bool {
		addresses = append(addresses, item.(string))
		return true
	})
	conn.connMgr.delAddr(addresses)

	// clean the address set's content
	conn.addrs.Clear()

	// The close operation will deliver websocket.WriteControl to send CloseMessage.
	// So the following close is not able to exec.
	// conn.close()
	//
	// should be responsible for closing the *websocket.Conn
	conn.Conn.CloseHandler()
}

func NewConnection(mgr *ConnMgr, conn *websocket.Conn, uid string) *Connection {
	return &Connection{
		connMgr: mgr,
		Conn:    conn,
		addrs:   set.New(set.ThreadSafe),
		Uid:     uid,
		Quit:    make(chan struct{}),
	}
}
