package controllers

import (
	"app/config"
	"app/models"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
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

type formLogin struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (a *Auth) Profile(c echo.Context) error {

	JwtVerify(c)

	sub := c.Get("sub")
	// var user models.User = sub.(models.User)

	return c.JSON(http.StatusOK, Map{"msg": sub})
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

func (a *Auth) SignIn(ctx echo.Context) error {

	var form formLogin
	if err := ctx.Bind(&form); err != nil {
		return ctx.JSON(http.StatusUnauthorized, Map{"status": "unable to bind data"})
	}

	var user models.User
	copier.Copy(&user, &form)

	if err := a.DB.Where("email = ?", form.Email).First(&user).Error; err != nil {
		return ctx.JSON(http.StatusUnauthorized, Map{"error": err.Error()})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
		return ctx.JSON(http.StatusUnauthorized, Map{"error": err.Error()})

	}

	// Create token
	at := jwt.New(jwt.SigningMethodHS256)
	// Set claims
	claims := at.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["name"] = user.Name
	claims["exp"] = time.Now().Add(time.Hour * 72).Local().Unix()

	// Generate encoded token and send it as response.
	token, err := at.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, Map{
		"token": token,
	})

}

func Authorize() echo.MiddlewareFunc {
	return middleware.JWT([]byte(os.Getenv("SECRET_KEY")))
}

func JwtVerify(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"]

	var users models.User

	db := config.GetDB()
	if err := db.First(&users, id).Error; err != nil {
		fmt.Println(err.Error())
	}
	c.Set("sub", users)

	return nil

}
