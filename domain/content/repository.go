package content

import (
	"contentgit/domain/content/projections"
	"contentgit/dtos"
	"context"
)

type ContentProjectionRepository interface {
	Create(ctx context.Context, projection projections.ContentProjection) error
	FindByID(ctx context.Context, tenantId string, id string) (*projections.ContentProjection, error)
	FindAll(ctx context.Context, tenantId string, pageable dtos.Pageable, sort *dtos.Sort) ([]projections.ContentProjection, int64, error)
	Save(ctx context.Context, projection *projections.ContentProjection) error
}
