package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/softstone1/woc/config"
	"github.com/softstone1/woc/infra/handler"

)

const (
	// server request timeout
	readHeaderTimeout = time.Second
	readTimeout       = 10 * time.Second
	handlerTimeout    = 40 * time.Second
	shutdownTimeout   = 5 * time.Second
)

type Mux struct {
	cfg            config.Env
	httpHandler    http.Handler
	weatherHandler *handler.Weather
}

// NewMux creates a new mux server and registers routes with the handlers.
// It also wraps the mux with logging and recovery middlewares
func NewMux(cfg config.Env, h *handler.Weather) (*Mux, error) {
	if h == nil {
		return nil, errors.New("handler is required")
	}
	mux := http.NewServeMux()
	// Register routes
	registerRoutes(mux, h)
	// Setup profiling routes
	if cfg.EnableProfiling() {
		setupProfiling(mux)
	}
	// wrap with logging and recovery middlewares
	wrappedMux := handlers.LoggingHandler(log.Writer(), handlers.RecoveryHandler()(mux))
	return &Mux{
		cfg:            cfg,
		httpHandler:    wrappedMux,
		weatherHandler: h,
	}, nil
}

// Run starts the server with graceful shutdown
func (s *Mux) Run(ctx context.Context) error {
	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%s", s.cfg.ServerPort()), // Call the s.cfg.ServerPort function
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		Handler:           http.TimeoutHandler(s.httpHandler, handlerTimeout, "request timed out"),
	}

	// Channel to communicate the server's closure
	done := make(chan error, 1) // Buffered channel to avoid goroutine leak
	// Start the server
	go func() {
		log.Printf("Starting server on %s", httpServer.Addr)
		// ListenAndServe always returns a non-nil error. After Shutdown or Close, the returned error is ErrServerClosed.
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			done <- err // Handle unexpected errors
		} else {
			done <- nil // No error, or graceful shutdown occurred.
		}
	}()

	// Wait for shutdown signal
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	log.Println("Shutting down server...")
	// Attempt to gracefully shut down the server
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	shutdownErr := httpServer.Shutdown(ctx)

	// Ensure that the done channel is also read to prevent goroutine leaks
	serverErr := <-done
	// Close the done channel since we're done receiving from it
	close(done)

	if shutdownErr != nil {
		return fmt.Errorf("error shutting down server: %w", shutdownErr)
	}
	if serverErr != nil {
		return fmt.Errorf("error running server: %w", serverErr)
	}
	return nil
}

// registerRoutes registers routes with the mux
func registerRoutes(mux *http.ServeMux, h *handler.Weather) {
	mux.HandleFunc("GET /", h.Home)
	mux.HandleFunc("GET /weather", h.GetWeatherByCity)
	mux.HandleFunc("GET /api/weather", h.GetWeatherByCityAPI)
}

func setupProfiling(mux *http.ServeMux) {
	slog.Info("🔍 Enable Profiling")
	mux.HandleFunc("GET /debug/pprof/", pprof.Index)
	mux.HandleFunc("GET /debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("GET /debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("GET /debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("GET /debug/pprof/trace", pprof.Trace)

	// Register other specific handlers
	mux.Handle("GET /debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("GET /debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("GET /debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	mux.Handle("GET /debug/pprof/block", pprof.Handler("block"))
}
