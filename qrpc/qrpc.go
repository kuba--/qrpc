package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/kuba--/qrpc"
)

var (
	cfg qrpc.Config
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.IntVar(&cfg.Port, "port", 9033, "port to listen on")
	flag.StringVar(&cfg.DataBasePath, "data", "/tmp/qrpc", "directory in which queue data is stored")
	flag.Uint64Var(&cfg.MaxCacheSize, "cache", 1024*1024, "max. cache size (bytes) before an item is evicted.")
	flag.DurationVar(&cfg.ClusterRequestTimeout, "timeout", 3*time.Second, "cluster request (gossip) timeout.")
	flag.DurationVar(&cfg.ClusterWatchInterval, "interval", time.Second, "cluster watch timer interval.")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [flags] [cluster peer(s) ...]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "flags:")
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()
	cfg.ClusterPeers = flag.Args()

	s := qrpc.NewServer(&cfg)
	errchan := s.Start()

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
