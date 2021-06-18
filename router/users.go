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

// GetUsersBookmarkHandler ブックマーク一覧を返す
func GetUsersBookmarkHandler(c echo.Context) error {
	user, err := GetMe(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Errorf("failed to get me: %w", err).Error())
	}

	bookmarks, err := model.GetBookmarks(c.Request().Context(), user.ID)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Errorf("failed to get bookmarks: %w", err).Error())
	}

	return c.JSON(http.StatusOK, bookmarks)
}
