package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/oinume/todomvc/backend/http_server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := http_server.New()
	router := server.NewRouter()
	log.Printf("Starting HTTP server on port %v", port)
	err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%v", port), router)
	if err != nil {
		log.Fatal(err)
	}
}
