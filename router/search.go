package router

import (
	"fmt"
	"net/http"

	"github.com/hackathon21spring-05/linq-backend/model"
	"github.com/labstack/echo/v4"
)

func GetSearchEntrys(c echo.Context) error {
	user, err := GetMe(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Errorf("failed to get me: %w", err).Error())
	}
	tag := c.QueryParam("tag")
	entryDetails, err := model.SearchEntrys(c.Request().Context(), tag, user.ID)

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Errorf("failed to search tag: %w", err).Error())
	}
	return c.JSON(http.StatusOK, entryDetails)
}
