package main

import (
	"context"
	"os"
	"os/signal"

	"wpld/internal/connectors/terminal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go handleSignals(cancel)

	terminal.Execute(ctx)
}

func handleSignals(cancel context.CancelFunc) {
	signCh := make(chan os.Signal)
	signal.Notify(signCh, os.Interrupt)

	for {
		sign := <-signCh
		switch sign {
		case os.Interrupt:
			cancel()
			return
		}
	}
}
