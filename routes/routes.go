package routes

import (
	"app/config"
	"app/controllers"
	"app/middlewares"

	"github.com/labstack/echo/v4"
)

func Serve(e *echo.Echo) {
	db := config.GetDB()
	v1 := e.Group("api/v1")
	// jwtVerify := controllers.JwtVerify()
	authenticate := middlewares.Authorize()

	authController := controllers.Auth{DB: db}
	authGroup := v1.Group("/auth")

	{
		authGroup.GET("", authController.Profile, authenticate)
		authGroup.POST("/sign-up", authController.SignUp)
		authGroup.POST("/sign-in", middlewares.SignIn)
	}

	productController := controllers.Products{DB: db}
	productGroup := v1.Group("/products")
	productGroup.Use(authenticate)
	{
		productGroup.GET("", productController.FindAll)
		productGroup.GET("/:id", productController.FindOne)
		productGroup.PUT("/:id", productController.Update)
		productGroup.DELETE("/:id", productController.Delete)
		productGroup.POST("", productController.Create)
	}
}
