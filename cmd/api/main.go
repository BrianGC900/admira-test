package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/admira-project/backend/internal/api"
	"github.com/admira-project/backend/internal/etl"
	"github.com/admira-project/backend/internal/storage"
	"github.com/admira-project/backend/internal/utils"
	"github.com/gorilla/mux"

	// "github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := configureLogger()

	httpClient := utils.NewRetryableHTTPClient(
		logger,
		getEnvAsInt("MAX_RETRIES", 3),
		getEnvAsInt("RETRY_BACKOFF_MS", 1000),
	)

	extractor := etl.NewExtractor(
		httpClient,
		os.Getenv("ADS_API_URL"),
		os.Getenv("CRM_API_URL"),
		logger,
	)

	transformer := etl.NewTransformer(logger)
	storage := storage.NewMemoryStorage()

	handler := api.NewHandler(extractor, transformer, storage, logger)

	router := mux.NewRouter()
	router.Use(loggingMiddleware(logger))

	router.HandleFunc("/ingest/run", handler.IngestHandler).Methods("POST")
	router.HandleFunc("/metrics/channel", handler.MetricsChannelHandler).Methods("GET")
	router.HandleFunc("/metrics/funnel", handler.MetricsFunnelHandler).Methods("GET")
	router.HandleFunc("/healthz", handler.HealthHandler).Methods("GET")
	router.HandleFunc("/readyz", handler.ReadyHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Infof("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server failed to start: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("Error during server shutdown: %v", err)
	}

	logger.Info("Server stopped")
}

func configureLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		level, err := logrus.ParseLevel(logLevel)
		if err == nil {
			logger.SetLevel(level)
		}
	}

	return logger
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

func loggingMiddleware(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := &responseWriter{w, http.StatusOK}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)

			logger.WithFields(logrus.Fields{
				"method":   r.Method,
				"path":     r.URL.Path,
				"status":   wrapped.status,
				"duration": duration.String(),
				"ip":       r.RemoteAddr,
			}).Info("HTTP request")
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
