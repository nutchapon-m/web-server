package mid

import (
	"context"
	"net/http"
	"time"

	"github.com/nutchapon-m/web-server/app/sdk/errs"
	"github.com/nutchapon-m/web-server/foundation/web"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(rate.Every(10*time.Second), 5)

func Limiter() web.MidFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx context.Context, r *http.Request) web.Encoder {
			if err := limiter.Wait(ctx); err != nil {
				return errs.Newf(errs.TooManyRequests, "Rate limit exceeded")
			}
			return next(ctx, r)
		}
	}
}
