package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	flag "github.com/spf13/pflag"
)

func main() {
	timeout := flag.DurationP("timeout", "t", 0, "Connection timeout")
	flag.Parse()
	if flag.NArg() < 2 {
		fmt.Printf("Usage: %s <host> <port>\n", os.Args[0])
		return
	}
	addr := net.JoinHostPort(flag.Arg(0), flag.Arg(1))
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	client := NewTelnetClient(addr, *timeout, os.Stdin, os.Stdout)
	defer client.Close()

	if err := client.ConnectContext(ctx); err != nil {
		log.Println(err)
		return
	}
	if err := client.Receive(); err != nil {
		log.Println(err)
		return
	}
	if err := client.Send(); err != nil {
		log.Println(err)
		return
	}
	<-ctx.Done()
	time.Sleep(1 * time.Second)
}
