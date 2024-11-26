package appservices

import (
	"contentgit/domain/content/commands"
	"contentgit/ports/out/persistance/eventsourcing"
)

type ContentService struct {
	Commands *commands.ContentCommands
}

func NewContentService(
	aggregateStore eventsourcing.AggregateStore,
) *ContentService {
	contentCommands := commands.NewContentCommands(
		commands.NewCreateUserSessionCmdHandler(aggregateStore),
		commands.NewUpdateContentFieldCmdHandler(aggregateStore),
		commands.NewAddContentFieldCommentCmdHandler(aggregateStore),
	)

	return &ContentService{Commands: contentCommands}
}
