package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/j7b/mailgun/webhook/types"
)

func ExampleDispatcher() {
	unsub := func(u *types.Unsubscribe) error {
		log.Println("Unsubscribed:", u.Recipient)
		return nil
	}
	d := &Dispatcher{
		Unsubscribe: unsub,
	}
	http.ListenAndServe(`127.0.0.1:0`, d)
}

func unsubscribes(chan *types.Unsubscribe) {}

func ExampleDispatcher_channels() {
	uchan := make(chan *types.Unsubscribe, 1024)
	go unsubscribes(uchan)
	unsub := func(u *types.Unsubscribe) error {
		select {
		case uchan <- u:
		default:
			return fmt.Errorf("rate limited")
		}
		return nil
	}
	d := &Dispatcher{
		Unsubscribe: unsub,
	}
	http.ListenAndServe(`127.0.0.1:0`, d)
}
