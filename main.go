package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/hackathon21spring-05/linq-backend/model"
	"github.com/hackathon21spring-05/linq-backend/router"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/go-sql-driver/mysql"
	sess "github.com/hackathon21spring-05/linq-backend/session"
)

var (
	db *sqlx.DB
)

func main() {
	log.Println("Server started")

	db, err := model.EstablishConnection()
	if err != nil {
		panic(err)
	}

	s, err := sess.NewSession(db.DB)
	if err != nil {
		panic(fmt.Errorf("failed in session constructor:%v", err))
	}

	e := echo.New()
	e.Use(session.Middleware(s.Store()))
	e.Use(middleware.Recover())

	clientID := os.Getenv("CLIENT_ID")
	if len(clientID) == 0 {
		panic(errors.New("ENV CLIENT_ID IS NULL"))
	}
	clientSecret := os.Getenv("CLIENT_SECRET")
	if len(clientSecret) == 0 {
		panic(errors.New("ENV CLIENT_SECRET IS NULL"))
	}

	err = router.SetRouting(e, s, clientID, clientSecret)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
