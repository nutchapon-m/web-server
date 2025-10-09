package mid

import (
	"context"
	"net/http"
	"path"

	"github.com/nutchapon-m/web-server/app/sdk/errs"
	"github.com/nutchapon-m/web-server/foundation/logger"
	"github.com/nutchapon-m/web-server/foundation/web"
)

// Errors handles errors coming out of the call chain.
func Errors(log *logger.Logger) web.MidFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx context.Context, r *http.Request) web.Encoder {
			resp := next(ctx, r)
			err := isError(resp)
			if err == nil {
				return resp
			}

			switch e := err.(type) {
			case *errs.Error:
				log.Error(ctx, "handled error during request",
					"err", err,
					"source_err_file", path.Base(e.FileName),
					"source_err_func", path.Base(e.FuncName))
				if e.Code == errs.InternalOnlyLog {
					e = errs.Newf(errs.Internal, "Internal Server Error")
				}
				return e
			case *errs.FieldErrors:
				log.Error(ctx, "handled error during request",
					"err", err,
					"source_err_file", path.Base(e.FileName),
					"source_err_func", path.Base(e.FuncName))
				return e
			default:
				return errs.Newf(errs.Internal, "Internal Server Error")
			}
		}
	}
}
