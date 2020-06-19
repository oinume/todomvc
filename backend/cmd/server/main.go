package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"contrib.go.opencensus.io/exporter/jaeger"
	_ "github.com/go-sql-driver/mysql"
	"go.opencensus.io/trace"
	"go.uber.org/zap"

	"github.com/oinume/todomvc/backend/config"
	"github.com/oinume/todomvc/backend/http_server"
	"github.com/oinume/todomvc/backend/logging"
)

func main() {
	logger, err := logging.New()
	if err != nil {
		log.Fatalf("logging.New failed: %v", err)
	}

	config.MustProcessDefault()

	// TODO: Stackdriver exporter
	exporter, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint:     config.DefaultVars.JaegerAgentEndpoint,
		CollectorEndpoint: config.DefaultVars.JaegerCollectorEndpoint,
		Process: jaeger.Process{
			ServiceName: "todomvc",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)

	db, err := sql.Open("mysql", config.DefaultVars.DBURL())
	// TODO: db.SetMaxIdleConns, etc...
	if err != nil {
		logger.Error("sql.Open failed", zap.Error(err))
		os.Exit(1)
	}

	server := http_server.New(db, logger)
	router := server.NewRouter()
	port := config.DefaultVars.HTTPPort
	logger.Info(fmt.Sprintf("Starting HTTP server on port %d", port))
	if err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%v", port), router); err != nil {
		logger.Fatal("http.ListenAndServe failed", zap.Error(err))
	}
	// TODO: graceful shutdown
}
