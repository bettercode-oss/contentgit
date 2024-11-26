package foundation

import (
	"contentgit/dtos"
	"sync"

	"gorm.io/gorm"
)

var (
	gormPaginatorOnce     sync.Once
	gormPaginatorInstance *gormPaginator
)

func GormPaginator() *gormPaginator {
	gormPaginatorOnce.Do(func() {
		gormPaginatorInstance = &gormPaginator{}
	})

	return gormPaginatorInstance
}

type gormPaginator struct {
}

func (gormPaginator) Pageable(pageable dtos.Pageable) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageable.Page > 0 {
			return db.Limit(pageable.PageSize).Offset(pageable.GetOffset())
		}
		return db
	}
}
