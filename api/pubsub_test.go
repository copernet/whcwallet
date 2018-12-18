package api

import (
	"context"
	"fmt"
	"github.com/copernet/whcwallet/model"
	"sync"
	"testing"
	"time"
)

func TestSubPub(t *testing.T) {
	channel := "wallet"

	c := make(chan []byte, 1)
	go func() {
		ctx, _ := context.WithCancel(context.Background())
		err := model.Subscribe(ctx, channel, c)
		if err != nil {
			t.Error("subscribe channel wallet, got error")
		}
	}()

	var wg sync.WaitGroup
	go func() {
		for {
			select {
			case ret := <-c:
				fmt.Println(string(ret))
				wg.Done()
			case <-time.After(2 * time.Second):
				fmt.Println("timeout")
			}
		}
	}()

	// to guarantee all subscriber have subscribed the specified channel
	time.Sleep(1 * time.Second)
	wg.Add(1)
	err := model.Publish(channel, "hello world")
	if err != nil {
		t.Error("publish message to channel wallet, got error")
	}

	wg.Wait()
}
