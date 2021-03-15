package main

import (
	"app/config"
	"app/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

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
	e.Static("/uploads", "./uploads")

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
