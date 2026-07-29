package main

import (
	"fmt"
	handler "go_text/api"
	"net/http"
)

func main() {
	http.HandleFunc("/api/health", handler.Handler)
	fmt.Println("listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
