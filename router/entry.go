package router

import (
	"net/http"

	"github.com/hackathon21spring-05/linq-backend/model"
	"github.com/labstack/echo/v4"
)

func GetEntryHandler(c echo.Context) error {
	// TODO
	return c.NoContent(http.StatusOK)
}

// PUT /entry ブックマークの追加・タグの更新・記事がなければ追加
// TODO: urlから記事のタイトルや中身を取得する
func PutEntryHandler(c echo.Context) error {
	req := struct {
		*model.Entry
		tags []string
	}{}

	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 記事がなければ追加
	err = model.AddEntry(c.Request().Context(), req.Entry)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// タグの追加（この機能いる？）

	// ブックマークに追加

	return c.NoContent(http.StatusOK)
}
