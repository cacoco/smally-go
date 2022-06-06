package main

import (
	"flag"

	h "github.com/cacoco/smally-go/pkg/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	httpPort = flag.String("http.port", ":8080", "Default HTTP port")
)

func main() {
	flag.Parse()

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/url", h.Create)
	e.GET("/:id", h.Get)

	// Start server
	e.Logger.Fatal(e.Start(*httpPort))
}
