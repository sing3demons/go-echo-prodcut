package controllers

import (
	"app/models"
	"net/http"

	"github.com/dgrijalva/jwt-go"
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
	Role  string `json:"role"`
}

func (a *Auth) CheckUser(c echo.Context) (*models.User, error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"]

	var users models.User

	if err := a.DB.First(&users, id).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

func (a *Auth) Profile(c echo.Context) error {
	user, _ := a.CheckUser(c)
	serializedUser := authResponse{}
	copier.Copy(&serializedUser, &user)
	return c.JSON(http.StatusOK, Map{"user": serializedUser})
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
