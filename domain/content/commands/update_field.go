package commands

import (
	"contentgit/domain/content"
	"contentgit/ports/out/persistance/eventsourcing"
	"context"
)

type UpdateContentField interface {
	Handle(ctx context.Context, cmd UpdateContentFieldCommand) error
}

type UpdateContentFieldCommand struct {
	AggregateID   string `json:"id"`
	TenantId      string `json:"tenantId"`
	FieldName     string `json:"fieldName"`
	BeforeValue   any    `json:"beforeValue"`
	AfterValue    any    `json:"afterValue"`
	CreatedById   string `json:"createdById"`
	CreatedByName string `json:"createdByName"`
}

type updateContentFieldCmdHandler struct {
	aggregateStore eventsourcing.AggregateStore
}

func (c *updateContentFieldCmdHandler) Handle(ctx context.Context, cmd UpdateContentFieldCommand) error {
	contentAggregate, err := content.NewContentAggregate(cmd.AggregateID, cmd.TenantId)
	if err != nil {
		return err
	}

	err = c.aggregateStore.Load(ctx, contentAggregate)
	if err != nil {
		return err
	}

	if err := contentAggregate.UpdateField(ctx, cmd.FieldName, cmd.BeforeValue, cmd.AfterValue, cmd.CreatedById, cmd.CreatedByName); err != nil {
		return err
	}

	return c.aggregateStore.Save(ctx, contentAggregate)
}

func NewUpdateContentFieldCmdHandler(aggregateStore eventsourcing.AggregateStore) *updateContentFieldCmdHandler {
	return &updateContentFieldCmdHandler{aggregateStore: aggregateStore}
}
