package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/oinume/todomvc/backend/infrastructure/mysql"

	"contrib.go.opencensus.io/exporter/jaeger"
	_ "github.com/go-sql-driver/mysql"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"go.uber.org/zap"

	"github.com/oinume/todomvc/backend/config"
	controller_http "github.com/oinume/todomvc/backend/controller/http"
	"github.com/oinume/todomvc/backend/logging"
)

func main() {
	logger, err := logging.New()
	if err != nil {
		log.Fatalf("logging.NewServer failed: %v", err)
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

	server := controller_http.NewServer(todoRepository, logger)
	router := server.NewRouter()
	port := config.DefaultVars.HTTPPort
	logger.Info(fmt.Sprintf("Starting HTTP server on port %d", port))
	ochttpHandler := &ochttp.Handler{
		Handler: router,
	}
	if err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%v", port), ochttpHandler); err != nil {
		logger.Fatal("http.ListenAndServe failed", zap.Error(err))
	}
	// TODO: graceful shutdown
}
