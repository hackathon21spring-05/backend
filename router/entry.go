package router

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/hackathon21spring-05/linq-backend/model"
	"github.com/labstack/echo/v4"
)

// Get GetEntryHandler新着を50件表示
func GetEntryHandler(c echo.Context) error {
	user, err := GetMe(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	entryDetails, err := model.GetNewEntrys(c.Request().Context(), user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, entryDetails)
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

	// URLパース
	u, err := url.Parse(req.Entry.Url)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to parse url: %w", err).Error())
	}
	removedUrl := u.Scheme + "://" + u.Host + u.Path
	if u.RawQuery != "" {
		removedUrl += "?" + u.RawQuery
	}
	usableURL, err := url.QueryUnescape(removedUrl)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to decode url: %w", err).Error())
	}
	entryId := model.ToHash(usableURL)

	// 記事の存在チェック
	count, err := model.FindEntry(c.Request().Context(), entryId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	// 新規記事のとき
	if count <= 0 {
		entry, err := getEntryContent(usableURL)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		// 記事の追加
		entry.ID = entryId
		err = model.AddEntry(c.Request().Context(), entry)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	// タグの追加
	err = model.AddTags(c.Request().Context(), entryId, req.Tags)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// ブックマークに追加
	user, err := GetMe(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

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
	user, err := GetMe(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	entryId := c.Param("entryId")
	entryDetail, err := model.GetEntryDetail(c.Request().Context(), user.ID, entryId)
	if err != nil {
		if err.Error() == "failed to get entry: sql: no rows in result set" {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, entryDetail)
}

func PostEntryTagHandler(c echo.Context) error {
	entryId := c.Param("entryId")
	tag := []string{c.Param("tag")}

	numEntrys, err := model.FindEntry(c.Request().Context(), entryId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if numEntrys == 0 {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	err = model.AddTags(c.Request().Context(), entryId, tag)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusCreated)
}

func DeleteEntryTagHandler(c echo.Context) error {

	entryId := c.Param("entryId")
	tag := []string{c.Param("tag")}

	err := model.DeleteTags(c.Request().Context(), entryId, tag)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func DeleteUsersBookmarkHandler(c echo.Context) error {

	entryId := c.Param("entryId")
	user, err := GetMe(c)

	err = model.DeleteBookmark(c.Request().Context(), entryId, user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
