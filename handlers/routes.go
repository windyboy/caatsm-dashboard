package handlers

import (
	"log/slog"
	"net/http"

	"casstm-dashboard/views"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	e.GET("/", HomeHandler)
}

func HomeHandler(ctx echo.Context) error {
	return renderView(ctx, views.Index())
}

func renderView(ctx echo.Context, cmp templ.Component) error {
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	if err := cmp.Render(ctx.Request().Context(), ctx.Response()); err != nil {
		slog.Error("Error rendering view", "err", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to render view")
	}
	return nil
}
