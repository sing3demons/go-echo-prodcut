package controllers

import (
	"math"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type pagingResult struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	PrevPage  int `json:"prevPage"`
	NextPage  int `json:"nextPage"`
	Count     int `json:"count"`
	TotalPage int `json:"totalPage"`
}
type pagination struct {
	ctx     echo.Context
	query   *gorm.DB
	records interface{}
}

func (p *pagination) pagingResource() *pagingResult {
	page, _ := strconv.Atoi(p.ctx.QueryParam("page"))
	if page == 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(p.ctx.QueryParam("limit"))
	if limit == 0 {
		limit = 12
	}

	ch := make(chan int)
	go p.countRecords(ch)

	offset := (page - 1) * limit
	p.query.Limit(limit).Offset(offset).Find(p.records)

	count := <-ch
	totalPage := int(math.Ceil(float64(count) / float64(limit)))

	var nextPage int
	if page == totalPage {
		nextPage = totalPage
	} else {
		nextPage = page + 1
	}

	return &pagingResult{
		Page:      page,
		Limit:     limit,
		PrevPage:  page - 1,
		NextPage:  nextPage,
		Count:     int(count),
		TotalPage: totalPage,
	}

}

func (p *pagination) countRecords(ch chan int) {
	var count int64
	p.query.Model(p.records).Count(&count)
	ch <- int(count)
}
