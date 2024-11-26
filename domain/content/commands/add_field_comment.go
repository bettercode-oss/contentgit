package commands

import (
	"contentgit/domain/content"
	"contentgit/ports/out/persistance/eventsourcing"
	"context"
)

type AddContentFieldComment interface {
	Handle(ctx context.Context, cmd AddContentFieldCommentCommand) error
}

type AddContentFieldCommentCommand struct {
	AggregateID   string `json:"id"`
	TenantId      string `json:"tenantId"`
	FieldName     string `json:"fieldName"`
	Comment       string `json:"comment"`
	CreatedById   string `json:"createdById"`
	CreatedByName string `json:"createdByName"`
}

type addContentFieldCommentCmdHandler struct {
	aggregateStore eventsourcing.AggregateStore
}

func (c *addContentFieldCommentCmdHandler) Handle(ctx context.Context, cmd AddContentFieldCommentCommand) error {
	contentAggregate, err := content.NewContentAggregate(cmd.AggregateID, cmd.TenantId)
	if err != nil {
		return err
	}

	err = c.aggregateStore.Load(ctx, contentAggregate)
	if err != nil {
		return err
	}

	if err := contentAggregate.AddFieldComment(ctx, cmd.FieldName, cmd.Comment, cmd.CreatedById, cmd.CreatedByName); err != nil {
		return err
	}

	return c.aggregateStore.Save(ctx, contentAggregate)
}

func NewAddContentFieldCommentCmdHandler(aggregateStore eventsourcing.AggregateStore) *addContentFieldCommentCmdHandler {
	return &addContentFieldCommentCmdHandler{aggregateStore: aggregateStore}
}
