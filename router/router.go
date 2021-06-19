package router

import (
	"net/http"
	"net/url"
	"strings"

	sess "github.com/hackathon21spring-05/linq-backend/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	s            sess.Session
	clientID     string
	clientSecret string
)

func SetRouting(e *echo.Echo, sess sess.Session, cltID string, cltSecret string) error {
	s = sess
	clientID = cltID
	clientSecret = cltSecret

	proxyConfig := middleware.DefaultProxyConfig
	clientURL, err := url.Parse("https://hackathon21spring-05.trap.show/linq-frontend/")
	if err != nil {
		panic(err)
	}
	proxyConfig.Balancer = middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
		{
			URL: clientURL,
		},
	})

	proxyConfig.Skipper = func(c echo.Context) bool {
		if strings.HasPrefix(c.Path(), "/api/") || strings.HasPrefix(c.Path(), "/openapi/") {
			return true
		}
		c.Request().Host = "hackathon21spring-05.trap.show"
		return false
	}
	proxyConfig.ModifyResponse = func(res *http.Response) error {
		res.Header.Set("Cache-Control", "max-age=3600")
		return nil
	}
	proxyConfig.Rewrite = map[string]string{
		"/callback": "/",
	}

	e.Use(middleware.ProxyWithConfig(proxyConfig))

	e.Static("/openapi", "docs/swagger")

	api := e.Group("/api")
	{
		// あとでセッション周りを適用する
		api.GET("/users/me", GetUsersMeHandler, UserAuthMiddleware)
		api.GET("/ping", pingHandler)
		apiEntry := api.Group("/entry")
		{
			apiEntry.GET("", GetEntryHandler)
			apiEntry.PUT("", PutEntryHandler)
			apiEntry.POST("/entry/:entryId/tag/:tag", PostEntryTagHandler)
			apiEntry.GET("/entry/:entryId", GetEntryDetailHandler)
			apiEntry.GET("/search", GetTagSearchHandler)
		}
		apiOAuth := api.Group("/oauth")
		{
			apiOAuth.GET("/callback", CallbackHandler)
			apiOAuth.POST("/generate/code", PostGenerateCodeHandler)
			apiOAuth.POST("/logout", PostGenerateCodeHandler, UserAuthMiddleware)
		}
	}

	return nil
}
