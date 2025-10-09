package mid

import (
	"context"
	"net/http"
	"time"

	"github.com/nutchapon-m/web-server/app/sdk/authclient"
	"github.com/nutchapon-m/web-server/app/sdk/errs"
	"github.com/nutchapon-m/web-server/foundation/web"
)

// Authenticate is a middleware function that integrates with an authentication client
// to validate user credentials and attach user data to the request context.
func Authenticate(client *authclient.Client) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			resp, err := client.Authenticate(ctx, r.Header.Get("authorization"))
			if err != nil {
				return errs.New(errs.Unauthenticated, err)
			}

			ctx = setUserID(ctx, resp.UserID)

			return next(ctx, r)
		}

		return h
	}

	return m
}
