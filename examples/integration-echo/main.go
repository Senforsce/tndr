package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/senforsce/tndr"
)

func main() {
	app := echo.New()
	app.GET("/", HomeHandler)
	app.Logger.Fatal(app.Start(":4000"))
}

// This custom Render replaces Echo's echo.Context.Render() with templ's templ.Component.Render().
func Render(ctx echo.Context, statusCode int, t tndr.Component) error {
	buf := tndr.GetBuffer()
	defer tndr.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

func HomeHandler(c echo.Context) error {
	return Render(c, http.StatusOK, Home())
}
