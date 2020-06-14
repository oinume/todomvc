package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"go.uber.org/zap"

	"github.com/oinume/todomvc/backend/http_server"
	"github.com/oinume/todomvc/backend/logging"
)

func main() {
	logger, err := logging.New()
	if err != nil {
		log.Fatalf("logging.New failed: %v", err)
	}

	server := http_server.New(logger)
	router := server.NewRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logger.Info("Starting HTTP server on port " + port)
	if err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%v", port), router); err != nil {
		logger.Fatal("http.ListenAndServe failed", zap.Error(err))
	}
	// TODO: graceful shutdown
}
