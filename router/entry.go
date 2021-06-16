package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetEntryHandler(c echo.Context) error {
	// TODO
	return c.NoContent(http.StatusOK)
}
