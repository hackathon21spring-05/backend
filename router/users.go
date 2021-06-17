package router

import (
	"fmt"
	"net/http"

	"github.com/hackathon21spring-05/linq-backend/model"
	"github.com/labstack/echo/v4"
)

// GetUsersMeHandler /users/meのハンドラー
func GetUsersMeHandler(c echo.Context) error {
	user, err := GetMe(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Errorf("failed to get me: %w", err).Error())
	}

	return c.JSON(http.StatusOK, user)
}

//GetUsersMeHandlerのまるこぴ
func GetUsersBookmarkHandler(c echo.Context) error {

	user, err := GetMe(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Errorf("failed to get me: %w", err).Error())
	}

	var userbookmarks model.UserBm

	userbookmarks, err = model.GetBookmark(c.Request().Context(), user)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Errorf("failed to get bookmark: %w", err).Error())
	}

	return c.JSON(http.StatusOK, userbookmarks)
}
