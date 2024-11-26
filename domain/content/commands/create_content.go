package commands

import (
	"contentgit/domain/content"
	"contentgit/ports/out/persistance/eventsourcing"
	"context"
)

type CreateContent interface {
	Handle(ctx context.Context, cmd CreateContentCommand) error
}

type CreateContentCommand struct {
	TenantID    string         `json:"tenantId"`
	AggregateID string         `json:"id"`
	Content     map[string]any `json:"content"`
}

type createContentCmdHandler struct {
	aggregateStore eventsourcing.AggregateStore
}

func (c *createContentCmdHandler) Handle(ctx context.Context, cmd CreateContentCommand) error {
	exists, err := c.aggregateStore.Exists(ctx, cmd.AggregateID)
	if err != nil {
		return err
	}
	if exists {
		return content.ErrContentAlreadyExists
	}

	contentAggregate, err := content.NewContentAggregate(cmd.AggregateID, cmd.TenantID)
	if err != nil {
		return err
	}

	err = contentAggregate.CreateContent(ctx, cmd.Content)
	if err != nil {
		return err
	}

	return c.aggregateStore.Save(ctx, contentAggregate)
}

func NewCreateUserSessionCmdHandler(aggregateStore eventsourcing.AggregateStore) *createContentCmdHandler {
	return &createContentCmdHandler{aggregateStore: aggregateStore}
}
