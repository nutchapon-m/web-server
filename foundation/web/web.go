package web

import (
	"context"
	"net/http"
	"strings"
)

var (
	allowMethods = []string{
		http.MethodGet,
		http.MethodOptions,
		http.MethodPost,
		http.MethodPatch,
		http.MethodPut,
		http.MethodDelete,
	}
	allowHeaders = []string{
		"Accept",
		"Content-Type",
		"Content-Length",
		"Accept-Encoding",
		"X-CSRF-Token",
		"Authorization",
	}
	exposeHeader = []string{}
)

type Encoder interface {
	Encode() (data []byte, contentType string, err error)
}

type HandlerFunc func(ctx context.Context, r *http.Request) Encoder

type Logger func(ctx context.Context, message string, args ...any)

type App struct {
	log     Logger
	mux     *http.ServeMux
	mw      []MidFunc
	origins []string
}

func NewApp(log Logger, mw ...MidFunc) *App {
	mux := http.NewServeMux()
	return &App{
		mux: mux,
		log: log,
		mw:  mw,
	}
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if a.origins != nil {
		reqOrigin := r.Header.Get("Origin")
		for _, origin := range a.origins {
			if origin == "*" || origin == reqOrigin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowMethods, ", "))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowHeaders, ", "))
		w.Header().Set("Access-Control-Expose-Headers", strings.Join(exposeHeader, ", "))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method != http.MethodOptions {
			w.Header().Set("Vary", "Origin")
		} else {
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Vary", "Access-Control-Request-Method")
			w.Header().Set("Vary", "Access-Control-Allow-Headers")
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

	a.mux.ServeHTTP(w, r)
}

func (a *App) EnableCORS(origins []string) {
	a.origins = origins
}

func (a *App) HandlerFunc(method, group, path string, handler HandlerFunc, mw ...MidFunc) {
	handler = wrapMiddleware(mw, handler)
	handler = wrapMiddleware(a.mw, handler)

	h := func(w http.ResponseWriter, r *http.Request) {
		// enforce HTTP method
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		ctx := setWriter(r.Context(), w)

		resp := handler(ctx, r)

		if err := Respond(ctx, w, resp); err != nil {
			a.log(ctx, "web-response")
			return
		}
	}

	pattern := path
	if group != "" {
		pattern = group + path
	}

	a.mux.HandleFunc(pattern, h)
}
