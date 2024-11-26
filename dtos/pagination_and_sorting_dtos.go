package dtos

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type PageResult[T any] struct {
	Result     T     `json:"result"`
	TotalCount int64 `json:"totalCount"`
}

const PageSize = 20

type Pageable struct {
	Page     int
	PageSize int
}

func (p Pageable) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

func NewPageableFromRequest(ctx *gin.Context) Pageable {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.Query("pageSize"))
	if err != nil {
		pageSize = PageSize
	}

	return Pageable{
		Page:     page,
		PageSize: pageSize,
	}
}

type Sort struct {
	Field     string
	Direction string
}

// NewSortFromRequest creates a new Sort object from the request query parameter sortBy.
// The sortBy query parameter should be in the format of "direction(field)".
// For example, "asc(id)" or "desc(name)".
// If the sortBy query parameter is not provided, it will return nil.
func NewSortFromRequest(ctx *gin.Context) *Sort {
	sortBy := ctx.Query("sortBy")
	if len(sortBy) == 0 {
		return nil
	}

	sortByParts := strings.Split(sortBy, "(")
	if len(sortByParts) != 2 {
		return nil
	}
	return &Sort{
		Field:     sortByParts[1][:len(sortByParts[1])-1],
		Direction: sortByParts[0],
	}
}
