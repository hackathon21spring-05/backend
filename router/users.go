package router

import (
	"fmt"
	"net/http"

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
	Bookmark, err := GetBm(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Errorf("failed to get bookmarks: %w", err).Error())
	}

	return c.JSON(http.StatusOK, Bookmark)
}