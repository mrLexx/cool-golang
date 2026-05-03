package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	timeout    time.Duration
	host, port string
)

type result struct {
	source string
	err    error
}

func init() {
	flag.DurationVar(&timeout, "timeout", time.Second*10, "connection timeout sec")
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.SetFlags(0)

	flag.Parse()
	if len(flag.CommandLine.Args()) < 2 {
		help()
		return
	}
	host = flag.CommandLine.Args()[0]
	port = flag.CommandLine.Args()[1]

	client := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Printf("...error: %v:", err)
		return
	}

	log.Printf("...Connected to %v:%v", host, port)

	results := make(chan result, 2)

	go func() {
		results <- result{"server", client.Receive()}
	}()

	go func() {
		results <- result{"client", client.Send()}
	}()

	select {
	case <-ctx.Done():
		log.Println("...Interrupted Ctrl+C")
	case r := <-results:
		switch r.source {
		case "server":
			if r.err != nil {
				log.Printf("...Connection lost by server: %v", r.err)
			} else {
				log.Println("...Connection was closed by peer")
			}
		case "client":
			if r.err != nil {
				log.Printf("...Client send error: %v", r.err)
			} else {
				log.Println("...EOF")
			}
		}
	}

	if err := client.Close(); err != nil {
		log.Printf("...Error when closing: %v:\n", err)
	}
	log.Printf("...The session closed")
}

func help() {
	log.Println("Usage: go-telnet [--timeout=NNs] host port")
	log.Println("          Command Summary:")
	log.Println("          --timeout=NNs - connection timeout (sec)")
	log.Println()
}
