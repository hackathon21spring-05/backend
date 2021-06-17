package router

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/hackathon21spring-05/linq-backend/model"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// GetMe 本番用のGetMe
func GetMe(c echo.Context) (*model.User, error) {
	sess, err := session.Get("sessions", c)
	if err != nil {
		return nil, fmt.Errorf("Failed In Getting Session:%w", err)
	}
	id := sess.Values["id"].(string)
	name := sess.Values["name"].(string)
	if len(id) == 0 || len(name) == 0 {
		accessToken := sess.Values["accessToken"].(string)
		if len(accessToken) == 0 {
			return nil, errors.New("AccessToken Is Null")
		}
		user, err := getMe(accessToken)
		if err != nil {
			return nil, fmt.Errorf("Failed In Getting Me:%w", err)
		}
		return user, nil
	}

	return &model.User{ID: id, Name: name}, nil
}

func GetBm(c echo.Context) error {
	
	//getme関数を呼び出してuser情報を取得
	user, err := GetMe(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Errorf("failed to get me: %w", err).Error())
	}

	//user_id=user.idとなる要素のentry_idをすべてrowsに代入
	rows, err := db.Query("SELECT entry_id FROM bookmarks where user_id=?", user.id)
	if err != nil {
		return c.NoContent(http.StatusNoContent)
	}
	
	//rowsをそのまま返す(これがうまく動作するかが未検証)
	return c.JSON(http.StatusOK, rows)
}

// UserAuthMiddleware 本番用のAPIにアクセスしたユーザーを認証するミドルウェア
func UserAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("sessions", c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("failed to get session:%w", err))
		}

		accessToken := sess.Values["accessToken"]
		if accessToken == nil {
			return c.NoContent(http.StatusUnauthorized)
		}
		c.Set("accessToken", accessToken)

		return next(c)
	}
}
