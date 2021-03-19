package controllers

import (
	"app/models"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)



type Products struct {
	DB *gorm.DB
}

type productForm struct {
	Name  string                `form:"name" validate:"required"`
	Desc  string                `form:"desc" validate:"required"`
	Price int                   `form:"price" validate:"required"`
	Image *multipart.FileHeader `form:"image" validate:"required"`
}

type productRespons struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Price int    `json:"price"`
	Image string `json:"image"`
}

type pagingRespons struct {
	Items  []productRespons `json:"items"`
	Paging *pagingResult    `json:"paging"`
}

//H - json formate
type H map[string]interface{}

func (p *Products) FindAll(ctx echo.Context) error {
	products := []models.Products{}

	pagination := pagination{
		ctx:     ctx,
		query:   p.DB,
		records: &products,
	}
	paging := pagination.pagingResource()
	// p.DB.Find(&products)

	serializedProducts := []productRespons{}
	copier.Copy(&serializedProducts, &products)
	return ctx.JSON(http.StatusOK, H{"products": pagingRespons{Items: serializedProducts, Paging: paging}})
}

//Create - inser product
func (p *Products) Create(ctx echo.Context) error {
	var form productForm
	if err := ctx.Bind(&form); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var product models.Products
	copier.Copy(&product, &form)

	if err := p.DB.Create(&product).Error; err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	p.setProductImage(ctx, &product)

	var serializedProduct productRespons
	copier.Copy(&serializedProduct, &product)
	return ctx.JSON(http.StatusOK, H{"product": serializedProduct})
}

//FindOne - find product
func (p *Products) FindOne(ctx echo.Context) error {
	product, err := p.findProductByID(ctx)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, H{"error": err.Error()})
	}

	var serializedProduct productRespons
	copier.Copy(&serializedProduct, &product)
	return ctx.JSON(http.StatusOK, H{"product": serializedProduct})
}

func (p *Products) Update(ctx echo.Context) error {
	var form productForm
	if err := ctx.Bind(&form); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, H{"error": err.Error()})
	}
	product, err := p.findProductByID(ctx)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, H{"error": err.Error()})
	}

	copier.Copy(&product, &form)

	if err := p.DB.Save(&product).Error; err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, H{"error": err.Error()})
	}

	p.setProductImage(ctx, product)

	var serializedProduct productRespons
	copier.Copy(&serializedProduct, &product)
	return ctx.JSON(http.StatusOK, H{"product": serializedProduct})
}

func (p *Products) Delete(ctx echo.Context) error {
	product, err := p.findProductByID(ctx)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, H{"error": err.Error()})
	}

	p.DB.Delete(&product)
	return ctx.NoContent(http.StatusOK)
}

func (p *Products) findProductByID(ctx echo.Context) (*models.Products, error) {
	var product models.Products
	id := ctx.Param("id")

	if err := p.DB.First(&product, id).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (p *Products) SaveFile(file *multipart.FileHeader, path string) error {
	// Source
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	// Destination
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	return nil
}

func (p *Products) setProductImage(ctx echo.Context, product *models.Products) error {
	file, err := ctx.FormFile("image")
	if err != nil || file == nil {
		return err
	}

	p.chekProductImage(ctx, product)

	path := "uploads/products/" + strconv.Itoa(int(product.ID))
	os.Mkdir(path, 0755)

	filename := path + "/" + file.Filename

	if err := p.SaveFile(file, filename); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, H{"error": err.Error()})
	}

	product.Image = os.Getenv("HOST") + "/" + filename

	if err := p.DB.Save(product).Error; err != nil {
		return err
	}

	return nil

}

func (p *Products) chekProductImage(ctx echo.Context, product *models.Products) error {
	if product.Image != "" {
		product.Image = strings.Replace(product.Image, os.Getenv("HOST"), "", 1)
		pwd, _ := os.Getwd()
		os.Remove(pwd + product.Image)
	}
	return nil
}
