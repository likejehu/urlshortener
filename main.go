package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/likejehu/urlshortener/cache"
	"github.com/likejehu/urlshortener/db"
	"github.com/likejehu/urlshortener/handlers"
	"github.com/likejehu/urlshortener/idgen"
)

func main() {
	//create a new instance of echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Routes
	handler := handlers.Handler{
		Store: db.LinksStore,
		Cache: cache.LinksCache,
		ID:    idgen.IDG,
	}

	e.POST("/createlink", handler.CreateLink)
	e.GET("/:shortUrl", handler.Redirect)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
