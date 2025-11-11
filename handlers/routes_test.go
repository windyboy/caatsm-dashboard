package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

type failingComponent struct{}

func (f failingComponent) Render(ctx context.Context, w io.Writer) error {
	return errors.New("render failure")
}

func TestRenderViewReturnsHTTPErrorOnFailure(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := renderView(ctx, failingComponent{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	httpErr, ok := err.(*echo.HTTPError)
	if !ok {
		t.Fatalf("expected *echo.HTTPError, got %T", err)
	}

	if httpErr.Code != http.StatusInternalServerError {
		t.Fatalf("expected HTTP 500, got %d", httpErr.Code)
	}
}
