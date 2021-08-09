package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sing3demons/go-echo-product/config"
	"github.com/sing3demons/go-echo-product/routes"

	"net/http"
)

type H map[string]interface{}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.InitDB()

	e := echo.New()
	e.HideBanner = true
	e.Static("/uploads", "./uploads")

	//middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}\n",
		Output: e.Logger.Output(),
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS()) // CORS default
	// e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{"https://labstack.com", "https://labstack.net"},
	// 	AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	// })) // CORS restricted

	e.GET("", homepage)
	routes.Serve(e)

	uploadDir := [...]string{"products", "users"}
	for _, dir := range uploadDir {
		os.MkdirAll("uploads/"+dir, 0755)
	}

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}

func homepage(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, H{"message": "product"})
}
