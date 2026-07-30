package main

import (
	"fmt"
	handler "go_text/api"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/api/health", handler.Handler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("listening on port %s...\n", port)
	http.ListenAndServe(":8080", nil)
}
