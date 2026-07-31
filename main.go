package main

import (
	"context"
	"fmt"
	"go_text/api"
	"go_text/internal/db"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err == nil {
		log.Println("[LOCAL] Loaded environment variables from .env")
	}
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = db.GetPool(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	log.Println("Database ready!")

	// 3. Serve HTTP on the established listener
	log.Fatal(http.Serve(listener, nil))
}
