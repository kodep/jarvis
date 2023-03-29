package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	onInterrupt(cancel)

	app, cleanup, err := InitializeApp()
	if err != nil {
		panic(err)
	}

	defer cleanup()

	app.Run(ctx)
}

func onInterrupt(fn func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		defer fn()
		<-c
	}()
}
