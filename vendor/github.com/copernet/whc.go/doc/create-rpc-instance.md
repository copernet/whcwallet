### RPC Instantiation :

All callers for `whc.go` SDK APIs are the Instantiated RPC object.

Example：

```Go
package rpc

import (
	"log"

	"github.com/copernet/whc.go/rpcclient"
)

var client *rpcclient.Client

func NewRPCInstance() *rpcclient.Client {
	if client != nil {
		return client
	}

	// Connect to local wormhole RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig {
		Host:         "39.104.152.31",
		User:         "5xX57X3Ad7KdjUv/lYFR",
		Pass:         "cz8S5tM/BUdQv6tyQdET",
		HTTPPostMode: true, 	// wormhole only supports HTTP POST mode
		DisableTLS:   true, 	// wormhole does not provide TLS by default
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
```

> The code of RPC instantiation is necessary in your own project. SDK only supplies a exported instantiation function: `rpcclient.New(connCfg, nil)`。

