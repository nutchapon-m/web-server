package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// NoResponse tells the Respond function to not respond to the request. In these
// cases the app layer code has already done so.
type NoResponse struct{}

// NewNoResponse constructs a no reponse value.
func NewNoResponse() NoResponse {
	return NoResponse{}
}

// Encode implements the Encoder interface.
func (NoResponse) Encode() ([]byte, string, error) {
	return nil, "", nil
}

// JSONResponse is a simple encoder for returning JSON payloads.
type JSONResponse struct {
	Status int
	Data   any
}

func JSON(status int, data any) JSONResponse {
	return JSONResponse{Status: status, Data: data}
}

func (j JSONResponse) HTTPStatus() int { return j.Status }

func (j JSONResponse) Encode() ([]byte, string, error) {
	if j.Data == nil {
		return nil, "application/json", nil
	}

	b, err := json.Marshal(j.Data)
	if err != nil {
		return nil, "", err
	}

	return b, "application/json; charset=utf-8", nil
}

// Response HTML
type HTMLResponse struct {
	Data string
}

func HTML(html string) HTMLResponse {
	return HTMLResponse{Data: html}
}

func (html HTMLResponse) HTTPStatus() int { return http.StatusOK }

func (html HTMLResponse) Encode() ([]byte, string, error) {
	if html.Data == "" {
		return nil, "text/html", nil
	}

	return []byte(html.Data), "text/html; charset=UTF-8", nil
}

// =====================================================================================================================

type httpStatus interface {
	HTTPStatus() int
}

// Respond sends a response to the client.
func Respond(ctx context.Context, w http.ResponseWriter, resp Encoder) error {
	if _, ok := resp.(NoResponse); ok {
		return nil
	}

	// If the context has been canceled, it means the client is no longer
	// waiting for a response.
	if err := ctx.Err(); err != nil {
		if errors.Is(err, context.Canceled) {
			return errors.New("client disconnected, do not send response")
		}
	}

	statusCode := http.StatusOK

	switch v := resp.(type) {
	case httpStatus:
		statusCode = v.HTTPStatus()

	case error:
		statusCode = http.StatusInternalServerError

	default:
		if resp == nil {
			statusCode = http.StatusNoContent
		}
	}

	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	data, contentType, err := resp.Encode()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("respond: encode: %w", err)
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)

	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("respond: write: %w", err)
	}

	return nil
}
