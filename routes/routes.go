package routes

import (
	"app/config"
	"app/controllers"

	"github.com/labstack/echo/v4"
)

func Serve(e *echo.Echo) {
	db := config.GetDB()
	v1 := e.Group("api/v1")

	authController := controllers.Auth{DB: db}
	authGroup := v1.Group("/auth")
	{
		authGroup.POST("/sign-up", authController.SignUp)
	}

	productController := controllers.Products{DB: db}
	productGroup := v1.Group("/products")
	{
		productGroup.GET("", productController.FindAll)
		productGroup.GET("/:id", productController.FindOne)
		productGroup.PUT("/:id", productController.Update)
		productGroup.DELETE("/:id", productController.Delete)
		productGroup.POST("", productController.Create)
	}
}
