package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/oinume/todomvc-example/backend/http_server"
)

func main() {
	server := http_server.New()
	router := server.NewRouter()
	err := http.ListenAndServe(fmt.Sprintf(":%d", 8080), router)
	if err != nil {
		log.Fatal(err)
	}
}
