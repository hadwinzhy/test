package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// PaginationParam is default query of pagination
type PaginationParam struct {
	Page    int    `form:"page,default=1" binding:"gt=0"`
	PerPage int    `form:"per_page,default=10" binding:"gt=0,max=50"`
	OrderBy string `form:"order_by,default=created_at"`
	SortBy  string `form:"sort_by,default=desc"`
}

// PaginationResponse will response of page
type PaginationResponse struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Per   int `json:"per"`
}

// SetPaginationToHeader will set pagination data to header
func SetPaginationToHeader(c *gin.Context, total int, page int, perPage int) {
	c.Writer.Header().Set("X-Total-Count", strconv.Itoa(total))
	c.Writer.Header().Set("X-Current-Page", strconv.Itoa(page))
	c.Writer.Header().Set("X-Per-Page", strconv.Itoa(perPage))
}

func SetPaginationToHeaderByStruct(c *gin.Context, pr PaginationResponse) {
	SetPaginationToHeader(c, pr.Total, pr.Page, pr.Per)
}
