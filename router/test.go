package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func pingHandler(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
