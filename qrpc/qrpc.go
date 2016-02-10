package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/kuba--/qrpc"
)

var (
	port  = flag.String("port", "2016", "port to listen on")
	data  = flag.String("data", ".qrpc", "directory in which queue data is stored")
	cache = flag.Uint64("cache", 1024*1024, "max cache size (bytes) before an item is evicted.")
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [flags] [cluster peer(s) ...]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "flags:")
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()
	peers := flag.Args()

	s := qrpc.NewServer(*data, *cache)
	errchan := s.Start(*port, peers...)

	// Handle SIGINT and SIGTERM.
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, os.Kill)

	select {
	case sig := <-sigchan:
		log.Printf("Signal %s received, shutting down...\n", sig)
		s.Stop()

	case err := <-errchan:
		log.Fatalln(err)
	}
}
