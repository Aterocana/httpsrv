package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Aterocana/httpsrv"
)

// must forces function to have no errors, otherwise it panics.
func must[T any](arg T, err error) T {
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(0x01)
	}
	return arg
}

func main() {
	opts := flags()
	srv := must(httpsrv.New(opts...))
	must[*uint8](nil, srv.Open())

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	<-sigc
	fmt.Println(" shutting down server")

	ctx, cancel := context.WithTimeoutCause(context.Background(), 5*time.Second, fmt.Errorf("timeout reached"))
	defer cancel()
	must[*uint8](nil, srv.Close(ctx))
}
