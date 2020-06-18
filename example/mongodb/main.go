package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/nelkinda/health-go"
	"github.com/nelkinda/health-go/checks/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	portPtr := flag.Int("port", 0, "Port for the backend service.")
	flag.Parse()

	listener, url := mustStart(*portPtr)
	_, _ = fmt.Fprintf(os.Stderr, "%s: info: URL: %s\n", os.Args[0], url)
	defer mustStop(listener)

	waitForIntOrTerm()
	os.Exit(0)
}

func waitForIntOrTerm() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
}

func mustStart(port int) (net.Listener, string) {
	url := "mongodb://127.0.0.1:27017"
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		panic(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		panic(err)
	}
	h := health.New(
		health.Health{
			Version:   "1",
			ReleaseID: "1.0.0-SNAPSHOT",
		},
		health.WithChecksProviders(
			mongodb.Health(url, client, time.Duration(1)*time.Second, time.Duration(200)*time.Millisecond),
		),
	)
	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.Handler)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		panic(err)
	}
	go func() {
		if err := http.Serve(listener, mux); err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				_, _ = fmt.Fprintf(os.Stderr, "Gracefully shutting down %v\n", listener.Addr())
			} else {
				panic(err)
			}
		}
	}()
	return listener, fmt.Sprintf("http://%v", listener.Addr())
}

func mustStop(closeable io.Closer) {
	if err := closeable.Close(); err != nil {
		panic(err)
	}
}
