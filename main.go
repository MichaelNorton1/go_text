package main

import (
	"fmt"
	handler "go_text/api"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", handler.Handler)
	http.HandleFunc("/api/health", handler.Handler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := "0.0.0.0:" + port
	fmt.Printf("listening on port %s...\n", port)
	http.ListenAndServe(":8080", nil)
	log.Fatal(http.ListenAndServe(addr, nil))
}
