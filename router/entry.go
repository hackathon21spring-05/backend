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
		Tags []string `json:"tags"`
	}{}

	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	entryId := model.ToHash(req.Entry.Url)

	// 記事がなければ追加
	err = model.AddEntry(c.Request().Context(), req.Entry)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// タグの追加
	err = model.AddTags(c.Request().Context(), entryId, req.Tags)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// ブックマークに追加
	user := model.User{}
	// user, err = GetMe(c)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	// }

	// session が使えないのでとりあえず仮に <TODO>
	user.ID = "060db77b-1d04-4686-a5ec-15c960159646"

	err = model.AddBookMark(c.Request().Context(), user.ID, entryId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 記事情報を取得
	var entryDetail model.EntryDetail
	entryDetail, err = model.GetEntryDetail(c.Request().Context(), user.ID, entryId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, entryDetail)
}

func GetEntryDetailHandler(c echo.Context) error {

	// ブックマークに追加
	user := model.User{}
	//user, err = GetMe(c)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	// }

	// session が使えないのでとりあえず仮に <TODO>
	user.ID = "060db77b-1d04-4686-a5ec-15c960159646"

	entryId := c.Param("entryId")

	numentrys, err := model.FindEntry(c.Request().Context(), entryId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if numentrys.Num == 0 {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	entrydetail, err := model.GetEntryDetail(c.Request().Context(), user.ID, entryId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, entrydetail)
}
