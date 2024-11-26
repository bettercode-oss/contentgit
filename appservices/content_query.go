package appservices

import (
	"contentgit/domain/content"
	"contentgit/domain/content/projections"
	"contentgit/dtos"
	"context"
)

type ContentQuery struct {
	contentProjectionRepository content.ContentProjectionRepository
}

func NewContentQuery(contentProjectionRepository content.ContentProjectionRepository) *ContentQuery {
	return &ContentQuery{contentProjectionRepository: contentProjectionRepository}
}

func (q ContentQuery) GetContents(context context.Context, tenantId string, pageable dtos.Pageable, sortable *dtos.Sort) ([]projections.ContentProjection, int64, error) {
	contents, totalCount, err := q.contentProjectionRepository.FindAll(context, tenantId, pageable, sortable)
	if err != nil {
		return nil, 0, err
	}
	return contents, totalCount, nil
}

func (q ContentQuery) GetContent(ctx context.Context, tenantId string, id string) (*projections.ContentProjection, error) {
	return q.contentProjectionRepository.FindByID(ctx, tenantId, id)
}
