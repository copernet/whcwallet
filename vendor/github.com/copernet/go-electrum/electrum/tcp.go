package electrum

import (
	"bufio"
	"crypto/tls"
	"log"
	"net"
	"time"
)

var DebugMode bool

const connTimeout = 3 * time.Second

type TCPTransport struct {
	conn      net.Conn
	responses chan []byte
	errors    chan error
}

func NewTCPTransport(addr string) (*TCPTransport, error) {
	conn, err := net.DialTimeout("tcp", addr, connTimeout)
	if err != nil {
		return nil, err
	}

	t := &TCPTransport{
		conn:      conn,
		responses: make(chan []byte),
		errors:    make(chan error),
	}
	go t.listen()

	return t, nil
}

func NewSSLTransport(addr string, config *tls.Config) (*TCPTransport, error) {
	dialer := net.Dialer{
		Timeout: connTimeout,
	}
	conn, err := tls.DialWithDialer(&dialer, "tcp", addr, config)
	if err != nil {
		return nil, err
	}

	t := &TCPTransport{
		conn:      conn,
		responses: make(chan []byte),
		errors:    make(chan error),
	}
	go t.listen()

	return t, nil
}

func (t *TCPTransport) SendMessage(body []byte) error {
	if DebugMode {
		log.Printf("%s [debug] %s <- %s", time.Now().Format("2006-01-02 15:04:05"), t.conn.RemoteAddr(), body)
	}

	_, err := t.conn.Write(body)
	return err
}

func (t *TCPTransport) listen() {
	defer t.conn.Close()
	reader := bufio.NewReader(t.conn)
	for {
		// The Node should send server.ping request continuously with a
		// reasonable break in order to keep connection alive. If not
		// client will receive a disconnection error and encounter
		// io.EOF with following os.Exit(1).
		line, err := reader.ReadBytes(delim)
		if err != nil {
			// block until start handle error
			t.errors <- err
			break
		}
		if DebugMode {
			log.Printf("%s [debug] %s -> %s", time.Now().Format("2006-01-02 15:04:05"), t.conn.RemoteAddr(), line)
		}

		t.responses <- line
	}
}

func (t *TCPTransport) Responses() <-chan []byte {
	return t.responses
}

func (t *TCPTransport) Errors() <-chan error {
	return t.errors
}
