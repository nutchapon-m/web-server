package mid

import (
	"context"
	"errors"

	"github.com/nutchapon-m/web-server/app/sdk/errs"
	"github.com/nutchapon-m/web-server/foundation/web"
)

// isError tests if the Encoder has an error inside of it.
func isError(e web.Encoder) error {
	err, isError := e.(error)
	if isError {
		var appErr *errs.Error
		if errors.As(err, &appErr) && appErr.Code == errs.None {
			return nil
		}
		return err
	}
	return nil
}

// =============================================================================

type ctxKey int

const (
	claimKey ctxKey = iota + 1
	userIDKey
	userKey
	productKey
	homeKey
	trKey
)

func setUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}
