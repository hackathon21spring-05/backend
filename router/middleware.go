package router

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/hackathon21spring-05/linq-backend/model"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// GetMe 本番用のGetMe
func GetMe(c echo.Context) (*model.User, error) {
	// OAuthのない開発環境では次の値をテストに用いる（後で消すと良いかも）
	env := os.Getenv("ENV")
	if env == "develop" {
		return &model.User{ID: "060db77b-1d04-4686-a5ec-15c960159646", Name: "toshi00"}, nil
	}

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

// UserAuthMiddleware 本番用のAPIにアクセスしたユーザーを認証するミドルウェア
func UserAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	// OAuthのない開発環境では次の値をテストに用いる（後で消すと良いかも）
	env := os.Getenv("ENV")
	if env == "develop" {
		return func(c echo.Context) error {
			return next(c)
		}
	}
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
