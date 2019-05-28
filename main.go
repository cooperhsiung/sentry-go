package main

import (
    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
    "gitlab.com/sentry-go/handler"
    "log"
)

// port 2385

func main() {

    e := echo.New()

    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    e.HideBanner = true

    e.POST("/point", handler.SaveHandler)
    e.GET("/health", func(c echo.Context) error {
        return c.String(200, "healthy~~")
    })

    log.Fatal(e.Start(":2385"))
}
