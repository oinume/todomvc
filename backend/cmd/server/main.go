package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"contrib.go.opencensus.io/exporter/jaeger"
	_ "github.com/go-sql-driver/mysql"
	"go.opencensus.io/trace"
	"go.uber.org/zap"

	"github.com/oinume/todomvc/backend/config"
	controller_http "github.com/oinume/todomvc/backend/controller/http"
	"github.com/oinume/todomvc/backend/infrastructure/mysql"
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
		Endpoint:          config.DefaultVars.JaegerAgentEndpoint,
		AgentEndpoint:     config.DefaultVars.JaegerAgentEndpoint,
		CollectorEndpoint: config.DefaultVars.JaegerCollectorEndpoint,
		Process: jaeger.Process{
			ServiceName: "todomvc",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	trace.RegisterExporter(exporter)

	db, err := mysql.NewDB(config.DefaultVars.DBURL())
	if err != nil {
		logger.Error("mysql.NewDB failed", zap.Error(err))
		os.Exit(1)
	}
	todoRepository := mysql.NewTodoRepository(db)

	addr := fmt.Sprintf("127.0.0.1:%v", config.DefaultVars.HTTPPort)
	server := controller_http.NewServer(addr, db, todoRepository, logger)
	logger.Info(fmt.Sprintf("Starting HTTP server on %s", addr))
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal("server.ListenAndServe failed", zap.Error(err))
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM)
	<-sigCh

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server.Shutdown failed", zap.Error(err))
	}
}
