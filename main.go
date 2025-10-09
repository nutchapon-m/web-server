package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/nutchapon-m/web-server/app/sdk/errs"
	"github.com/nutchapon-m/web-server/app/sdk/mux"
	"github.com/nutchapon-m/web-server/foundation/env"
	"github.com/nutchapon-m/web-server/foundation/logger"
	"github.com/nutchapon-m/web-server/foundation/web"
)

var (
	build = flag.String("mode", "develop", "Service running on mode: develop or release")
	port  = flag.String("port", "8000", "Service port")
)

func main() {
	// -------------------------------------------------------------------------
	// Flag on start service

	flag.Parse()

	// -------------------------------------------------------------------------
	// Load ENV

	env.Load(".", "config.yml")

	// -------------------------------------------------------------------------
	// Start service

	log := logger.New(os.Stdout, logger.LevelInfo, "WEB-API")
	ctx := context.Background()
	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	// -------------------------------------------------------------------------
	// GOMAXPROCS

	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------
	// Configuration

	cfg := mux.Config{
		Build: *build,
		Log:   log,
	}

	var addr string
	if *build == "develop" {
		addr = fmt.Sprintf("localhost:%s", *port)
	} else {
		addr = fmt.Sprintf(":%s", *port)
	}

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

	if err := server.Shutdown(ctx); err != nil {
		log.Error(ctx, "shutdown error", "err", err)
	}

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

	app.HandlerFunc(http.MethodGet, "/api", "/test-error", func(ctx context.Context, r *http.Request) web.Encoder {
		return errs.NewFieldErrors("value", errors.New("field value is required"))
	})
}
