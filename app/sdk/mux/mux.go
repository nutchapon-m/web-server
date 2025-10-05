package mux

import (
	"net/http"

	"github.com/nutchapon-m/web-server/app/sdk/mid"
	"github.com/nutchapon-m/web-server/foundation/logger"
	"github.com/nutchapon-m/web-server/foundation/web"
)

// Options represent optional parameters.
type Options struct {
	corsOrigin []string
}

// WithCORS provides configuration options for CORS.
func WithCORS(origins []string) func(opts *Options) {
	return func(opts *Options) {
		opts.corsOrigin = origins
	}
}

type Config struct {
	Build string
	Log   *logger.Logger
}

type RouteAdder interface {
	Add(app *web.App, cfg Config)
}

func WebAPI(cfg Config, routeAdder RouteAdder, options ...func(opts *Options)) http.Handler {
	app := web.NewApp(
		cfg.Log.Info,
		mid.Logger(cfg.Log),
		mid.Errors(cfg.Log),
		mid.Panics(),
		mid.CSRF(),
	)

	var opts Options
	for _, option := range options {
		option(&opts)
	}

	if len(opts.corsOrigin) > 0 {
		app.EnableCORS(opts.corsOrigin)
	}

	routeAdder.Add(app, cfg)
	return app
}
