package routes

import (
	"github.com/sing3demons/go-echo-product/config"
	"github.com/sing3demons/go-echo-product/controllers"
	"github.com/sing3demons/go-echo-product/middlewares"

	"github.com/labstack/echo/v4"
)

func Serve(e *echo.Echo) {
	db := config.GetDB()
	cache := config.NewRedisCache("redis:6379", 1, 10)
	v1 := e.Group("api/v1")
	// jwtVerify := controllers.JwtVerify()
	authenticate := middlewares.Authorize()

	authController := controllers.Auth{DB: db}
	authGroup := v1.Group("/auth")

	{
		authGroup.GET("/profile", authController.Profile, authenticate)
		authGroup.POST("/sign-up", authController.SignUp)
		authGroup.POST("/sign-in", middlewares.SignIn)
	}

	productController := controllers.Products{DB: db, Cache: cache}
	productGroup := v1.Group("/products")

	// productGroup.Use(authenticate)
	{
		productGroup.GET("", productController.FindAll)
		productGroup.GET("/:id", productController.FindOne, authenticate)
		productGroup.PUT("/:id", productController.Update, authenticate)
		productGroup.DELETE("/:id", productController.Delete, authenticate)
		productGroup.POST("", productController.Create, authenticate)
	}
}
