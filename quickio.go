package main

import (
	"fmt"
	//"os"
	"github.com/iheartradio/quickigo"
)

func main() {
	qio := quickigo.New("ws://quickio.iheart.com")

	ch := make(chan quickigo.Cb)
	opened := make(chan quickigo.Cb)
	closed := make(chan quickigo.Cb)
	error := make(chan quickigo.Cb)

	qio.On(quickigo.EvOpen, nil, opened, nil)
	qio.On(quickigo.EvClose, nil, closed, nil)
	qio.On(quickigo.EvError, nil, error, nil)
	qio.Open()

	for i := 0; i < 2; i++ {
		select {
		case <-opened:
			fmt.Println("connection established")
			qio.Close()

		case <-closed:
			fmt.Println("lost connection")

		case cb := <-error:
			fmt.Println("connection error:", cb.Err)
		}
	}

	qio.Open()
	<-opened

	ctx := "some data"
	qio.Send("/clienttest/echo", "echo data", ch, ctx)

	cb := <-ch
	fmt.Println(cb.Data, ":", cb.Ctx)

	// Check if the server wanted a callback
	if cb.CanReply() {
		cb.Reply("some more data", ch, ctx)
		cb = <-ch
		fmt.Println(cb.Data)
	}

	qio.Close()
	// At this point, it's completely safe to call On()/Off()/Send();
	// any relevant state will be sent to the server when a connection
	// is reestablished.
}
