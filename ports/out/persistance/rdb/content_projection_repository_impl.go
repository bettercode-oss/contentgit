package rdb

import (
	"contentgit/domain/content/projections"
	"contentgit/dtos"
	"contentgit/foundation"
	persistence "contentgit/ports/out/persistance"
	"context"
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type ContentProjectionRepositoryImpl struct {
}

func (ContentProjectionRepositoryImpl) Create(ctx context.Context, projection projections.ContentProjection) error {
	db := foundation.ContextProvider().GetDB(ctx)

	if err := db.Create(&projection).Error; err != nil {
		return errors.Wrap(err, "db error")
	}

	return nil
}

func (ContentProjectionRepositoryImpl) FindByID(ctx context.Context, tenantId string, id string) (*projections.ContentProjection, error) {
	db := foundation.ContextProvider().GetDB(ctx)

	var projection projections.ContentProjection
	if err := db.Preload("FieldChanges").Preload("FieldComments").First(&projection, "tenant_id = ? AND id = ?", tenantId, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, persistence.ErrRecordNotFound
		}

		return nil, errors.Wrap(err, "db error")
	}

	return &projection, nil
}

func (ContentProjectionRepositoryImpl) FindAll(ctx context.Context, tenantId string, pageable dtos.Pageable, sort *dtos.Sort) ([]projections.ContentProjection, int64, error) {
	db := foundation.ContextProvider().GetDB(ctx).Model(&projections.ContentProjection{})

	var entities = make([]projections.ContentProjection, 0)
	var totalCount int64

	db = db.Where("tenant_id = ?", tenantId)

	if sort != nil {
		db = db.Order(fmt.Sprintf("%s %s", sort.Field, sort.Direction))
	}

	if err := db.Count(&totalCount).Scopes(foundation.GormPaginator().Pageable(pageable)).Find(&entities).Error; err != nil {
		return entities, totalCount, errors.Wrap(err, "db error")
	}

	return entities, totalCount, nil
}

func (ContentProjectionRepositoryImpl) Save(ctx context.Context, entity *projections.ContentProjection) error {
	db := foundation.ContextProvider().GetDB(ctx)

	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Save(entity).Error; err != nil {
		return errors.Wrap(err, "db error")
	}

	return nil
}
