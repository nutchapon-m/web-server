package mid

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	"github.com/nutchapon-m/web-server/app/sdk/errs"
	"github.com/nutchapon-m/web-server/foundation/web"
)

var (
	csrfTokenKey    = "X-CSRF-Token"
	csrfAllowMethod = []string{"GET", "HEAD", "OPTIONS", "TRACE"}

	storer map[string]string // build-in memory cache in future.
)

func CSRF() web.MidFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx context.Context, r *http.Request) web.Encoder {

			if slices.Contains(csrfAllowMethod, r.Method) {
				return next(ctx, r)
			}

			val := r.Header.Get(csrfTokenKey)

			var appErr *errs.Error
			if val == "" {
				return errs.Newf(errs.PermissionDenied, "The csrf token is required")
			}

			key := fmt.Sprintf("ua:csrftoken:%s", val)
			if _, exsits := storer[key]; exsits {
				return next(ctx, r)
			}

			appErr = errs.Newf(errs.PermissionDenied, "Invalid csrf token")
			return appErr
		}
	}
}
