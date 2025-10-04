package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nutchapon-m/web-server/app/sdk/mux"
	"github.com/nutchapon-m/web-server/foundation/logger"
	"github.com/nutchapon-m/web-server/foundation/web"
)

func main() {
	log := logger.New(os.Stdout, logger.LevelInfo, "TEST")
	ctx := context.Background()
	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	cfg := mux.Config{
		Log: log,
	}

	addr := "localhost:8000"
	server := http.Server{
		Addr:         addr,
		Handler:      mux.WebAPI(cfg, buildRoutes()),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	log.Info(ctx, "Server running", "addr", addr)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// log the error and trigger shutdown
			log.Error(ctx, "listen and serve", "err", err)
			shutdown <- os.Interrupt
		}
	}()

	<-shutdown
	log.Info(ctx, "shutdown")
	return nil
}

func buildRoutes() mux.RouteAdder {
	return Routes()
}

// =====================================================================================================================

func Routes() add {
	return add{}
}

type add struct{}

func (add) Add(app *web.App, cfg mux.Config) {
	app.HandlerFunc(http.MethodGet, "/api", "/test", func(ctx context.Context, r *http.Request) web.Encoder {
		return web.JSON(http.StatusOK, map[string]any{"ok": true, "path": r.URL.Path})
	})

	app.HandlerFunc(http.MethodPost, "/api", "/user", func(ctx context.Context, r *http.Request) web.Encoder {
		return web.JSON(http.StatusOK, map[string]any{"ok": true, "path": r.URL.Path})
	})
}
