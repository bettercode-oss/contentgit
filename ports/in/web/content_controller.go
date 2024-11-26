package web

import (
	"contentgit/app/datasource"
	"contentgit/appservices"
	"contentgit/domain/content"
	"contentgit/domain/content/commands"
	"contentgit/domain/content/projections"
	"contentgit/dtos"
	"contentgit/foundation"
	persistence "contentgit/ports/out/persistance"
	"context"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type ContentController struct {
	routerGroup    *gin.RouterGroup
	contentService *appservices.ContentService
	contentQuery   *appservices.ContentQuery
}

func NewContentController(rg *gin.RouterGroup, contentService *appservices.ContentService, contentQuery *appservices.ContentQuery) *ContentController {
	return &ContentController{
		routerGroup:    rg,
		contentService: contentService,
		contentQuery:   contentQuery,
	}
}

func (controller ContentController) MapRoutes() {
	route := controller.routerGroup.Group("/tenants/:tenantId/contents")
	route.POST("bulk", controller.createBulkContents)
	route.POST("", controller.createContent)
	route.GET("", controller.getContents)
	route.GET(":id", controller.getContent)
	route.PUT(":id/:fieldName", controller.updateContentField)
	route.POST(":id/:fieldName/comments", controller.addFieldComment)
}

func (controller ContentController) createBulkContents(ctx *gin.Context) {
	tenantId := ctx.Param("tenantId")
	if len(tenantId) == 0 {
		ctx.JSON(http.StatusBadRequest, "tenantId is required")
		return
	}

	var contents []any

	if err := ctx.BindJSON(&contents); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := datasource.TransactionalWithContext(ctx.Request.Context(), func(ctx context.Context) error {
		for _, c := range contents {
			command := commands.CreateContentCommand{
				TenantID:    tenantId,
				AggregateID: uuid.New().String(),
				Content:     c.(map[string]any),
			}

			err := controller.contentService.Commands.CreateContent.Handle(ctx, command)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		foundation.GinErrorHandler().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (controller ContentController) createContent(ctx *gin.Context) {
	tenantId := ctx.Param("tenantId")
	if len(tenantId) == 0 {
		ctx.JSON(http.StatusBadRequest, "tenantId is required")
		return
	}

	var content any

	if err := ctx.BindJSON(&content); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := datasource.TransactionalWithContext(ctx.Request.Context(), func(ctx context.Context) error {
		command := commands.CreateContentCommand{
			TenantID:    tenantId,
			AggregateID: uuid.New().String(),
			Content:     content.(map[string]any),
		}

		return controller.contentService.Commands.CreateContent.Handle(ctx, command)
	})

	if err != nil {
		foundation.GinErrorHandler().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (controller ContentController) getContent(ctx *gin.Context) {
	tenantId := ctx.Param("tenantId")
	if len(tenantId) == 0 {
		ctx.JSON(http.StatusBadRequest, "tenantId is required")
		return
	}

	id := ctx.Param("id")
	if len(id) == 0 {
		ctx.JSON(http.StatusBadRequest, "id is required")
		return
	}

	contentProjection, err := controller.contentQuery.GetContent(ctx.Request.Context(), tenantId, id)
	if err != nil {
		if errors.Is(err, persistence.ErrRecordNotFound) {
			ctx.Status(http.StatusNotFound)
			return
		}
		foundation.GinErrorHandler().InternalServerError(ctx, err)
		return
	}

	contentDetails := dtos.ContentDetails{
		Id:        contentProjection.Id,
		Content:   contentProjection.Content,
		CreatedAt: contentProjection.CreatedAt,
		UpdatedAt: contentProjection.UpdatedAt,
	}

	contentDetails.FieldComments = make([]dtos.ContentDetailsFieldComment, 0)
	fieldCommentsGroupedByFieldName := foundation.GroupByProperty(contentProjection.FieldComments, func(fieldComment projections.ContentFieldComment) string {
		return fieldComment.Name
	})
	for fieldName, group := range fieldCommentsGroupedByFieldName {
		fieldComment := dtos.ContentDetailsFieldComment{
			Field: fieldName,
		}

		comments := make([]dtos.ContentDetailsComment, 0)
		for _, comment := range group {
			comments = append(comments, dtos.ContentDetailsComment{
				Id:            comment.ID,
				Comment:       comment.Comment,
				CreatedAt:     comment.CreatedAt,
				CreatedById:   comment.CreatedById,
				CreatedByName: comment.CreatedByName,
			})
		}
		sort.Slice(comments, func(i, j int) bool {
			return comments[i].Id < comments[j].Id
		})
		fieldComment.Comments = comments
		contentDetails.FieldComments = append(contentDetails.FieldComments, fieldComment)
	}

	contentDetails.FieldChanges = make([]dtos.ContentDetailsFieldChange, 0)
	fieldChangesGroupedByFieldName := foundation.GroupByProperty(contentProjection.FieldChanges, func(fieldChange projections.ContentFieldChange) string {
		return fieldChange.Name
	})
	for fieldName, group := range fieldChangesGroupedByFieldName {
		fieldChange := dtos.ContentDetailsFieldChange{
			Field: fieldName,
		}

		changes := make([]dtos.ContentDetailsUpdateField, 0)
		for _, change := range group {
			changes = append(changes, dtos.ContentDetailsUpdateField{
				Id:            change.ID,
				BeforeValue:   change.Content.BeforeValue,
				AfterValue:    change.Content.AfterValue,
				CreatedAt:     change.CreatedAt,
				CreatedById:   change.Content.CreatedById,
				CreatedByName: change.Content.CreatedByName,
			})
		}
		sort.Slice(changes, func(i, j int) bool {
			return changes[i].Id > changes[j].Id
		})
		fieldChange.Changes = changes
		contentDetails.FieldChanges = append(contentDetails.FieldChanges, fieldChange)
	}

	ctx.JSON(http.StatusOK, contentDetails)
}

func (controller ContentController) getContents(ctx *gin.Context) {
	tenantId := ctx.Param("tenantId")
	if len(tenantId) == 0 {
		ctx.JSON(http.StatusBadRequest, "tenantId is required")
		return
	}

	pageable := dtos.NewPageableFromRequest(ctx)
	sortable := dtos.NewSortFromRequest(ctx)

	contentProjections, totalCount, err := controller.contentQuery.GetContents(ctx.Request.Context(), tenantId, pageable, sortable)
	if err != nil {
		foundation.GinErrorHandler().InternalServerError(ctx, err)
		return
	}

	contentSummaries := make([]dtos.ContentSummary, 0)
	for _, contentEntity := range contentProjections {
		contentSummaries = append(contentSummaries, dtos.ContentSummary{
			Id:        contentEntity.Id,
			Content:   contentEntity.Content,
			CreatedAt: contentEntity.CreatedAt,
			UpdatedAt: contentEntity.UpdatedAt,
		})
	}

	pageResult := dtos.PageResult[[]dtos.ContentSummary]{
		Result:     contentSummaries,
		TotalCount: totalCount,
	}

	ctx.JSON(http.StatusOK, pageResult)
}

func (controller ContentController) updateContentField(ctx *gin.Context) {
	tenantId := ctx.Param("tenantId")
	if len(tenantId) == 0 {
		ctx.JSON(http.StatusBadRequest, "tenantId is required")
		return
	}

	id := ctx.Param("id")
	fieldName := ctx.Param("fieldName")
	if len(id) == 0 || len(fieldName) == 0 {
		ctx.JSON(http.StatusBadRequest, "id and fieldName are required")
		return
	}

	var updateField dtos.ContentUpdateField
	if err := ctx.BindJSON(&updateField); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := datasource.TransactionalWithContext(ctx.Request.Context(), func(ctx context.Context) error {
		command := commands.UpdateContentFieldCommand{
			AggregateID:   id,
			TenantId:      tenantId,
			FieldName:     fieldName,
			BeforeValue:   updateField.BeforeValue,
			AfterValue:    updateField.AfterValue,
			CreatedById:   updateField.CreatedById,
			CreatedByName: updateField.CreatedByName,
		}

		return controller.contentService.Commands.UpdateContentField.Handle(ctx, command)
	})

	if err != nil {
		if errors.Is(err, persistence.ErrRecordNotFound) || errors.Is(err, content.ErrFieldNotFound) {
			ctx.Status(http.StatusNotFound)
			return
		}

		if errors.Is(err, content.ErrFieldUpdateConflict) {
			ctx.Status(http.StatusConflict)
			return
		}

		foundation.GinErrorHandler().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (controller ContentController) addFieldComment(ctx *gin.Context) {
	tenantId := ctx.Param("tenantId")
	if len(tenantId) == 0 {
		ctx.JSON(http.StatusBadRequest, "tenantId is required")
		return
	}

	id := ctx.Param("id")
	fieldName := ctx.Param("fieldName")
	if len(id) == 0 || len(fieldName) == 0 {
		ctx.JSON(http.StatusBadRequest, "id and fieldName are required")
		return
	}

	var fieldComment dtos.ContentFieldComment
	if err := ctx.BindJSON(&fieldComment); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := datasource.TransactionalWithContext(ctx.Request.Context(), func(ctx context.Context) error {
		command := commands.AddContentFieldCommentCommand{
			AggregateID:   id,
			TenantId:      tenantId,
			FieldName:     fieldName,
			Comment:       fieldComment.Comment,
			CreatedById:   fieldComment.CreatedById,
			CreatedByName: fieldComment.CreatedByName,
		}

		return controller.contentService.Commands.AddContentFieldComment.Handle(ctx, command)
	})

	if err != nil {
		if errors.Is(err, persistence.ErrRecordNotFound) || errors.Is(err, content.ErrFieldNotFound) {
			ctx.Status(http.StatusNotFound)
			return
		}

		foundation.GinErrorHandler().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
