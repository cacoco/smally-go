package main

import (
	h "github.com/cacoco/smally-go/pkg/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/url", h.Create)
	e.GET("/:id", h.Get)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
