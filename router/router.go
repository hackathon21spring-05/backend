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

	api := e.Group("")
	api.Use(UserAuthMiddleware)

	apiOAuth := api.Group("/oauth")
	{
		apiOAuth.GET("/callback", CallbackHandler)
		apiOAuth.GET("/generate/code", PostGenerateCodeHandler)
		apiOAuth.POST("/logout", PostGenerateCodeHandler)
	}
	return nil
}
