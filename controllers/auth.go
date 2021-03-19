package controllers

import (
	"app/models"
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Map map[string]interface{}

type Auth struct {
	DB *gorm.DB
}

type authForm struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=4"`
	Name     string `json:"name" validate:"required"`
}

type authResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name" `
}

func (a *Auth) SignUp(ctx echo.Context) error {

	var form authForm
	if err := ctx.Bind(&form); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, Map{"error": err.Error()})
	}

	var user models.User

	copier.Copy(&user, &form)
	user.Password = user.GenerateEncryptedPassword()

	if err := a.DB.Create(&user).Error; err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, Map{"error": err.Error(), "message": "ลงทะเบียนไม่สำเร็จ"})

	}

	var serializedUser authResponse
	copier.Copy(&serializedUser, &user)

	return ctx.JSON(http.StatusCreated, Map{"user": serializedUser, "message": "ลงทะเบียนสำเร็จ"})
}
