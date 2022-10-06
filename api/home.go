package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Home(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to NFTime")
}

func RedirectTest(c echo.Context) error {
	return c.Redirect(http.StatusFound, "/")
}
