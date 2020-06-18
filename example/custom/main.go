package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/nelkinda/health-go"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

type custom struct{}

func (*custom) HealthChecks(ctx context.Context) map[string][]health.Checks {
	return map[string][]health.Checks{"custom": {{ComponentID: "custom-component", Status: health.Pass}}}
}

func (*custom) AuthorizeHealth(r *http.Request) bool {
	return true
}

func mustStart(port int) (net.Listener, string) {
	h := health.New(health.Health{Version: "1", ReleaseID: "1.0.0-SNAPSHOT"}, health.WithChecksProviders(&custom{}))
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
