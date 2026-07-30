package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"go_text/api"
)

func main() {
	// Register handlers for root and API
	http.HandleFunc("/", handler.Handler)
	http.HandleFunc("/api/health", handler.Handler)
	http.HandleFunc("/api/notify", handler.NotifyHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 1. Force explicit IPv4 binding FIRST
	listener, err := net.Listen("tcp4", "0.0.0.0:"+port)
	if err != nil {
		log.Fatalf("Failed to bind port %s: %v", port, err)
	}

	// 2. Only print log AFTER the TCP socket is confirmed open
	fmt.Printf("Server actively listening on %s...\n", listener.Addr().String())

	// 3. Serve HTTP on the established listener
	log.Fatal(http.Serve(listener, nil))
}
