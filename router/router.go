package router

import (
	sess "github.com/hackathon21spring-05/linq-backend/session"
	"github.com/labstack/echo/v4"
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

	api := e.Group("/api")
	{
		apiEntry := api.Group("/entry")
		{
			apiEntry.GET("/", GetEntryHandler)
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
