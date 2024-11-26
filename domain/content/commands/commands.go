package commands

type ContentCommands struct {
	CreateContent
	UpdateContentField
	AddContentFieldComment
}

func NewContentCommands(
	createContent CreateContent,
	updateContentField UpdateContentField,
	addContentFieldComment AddContentFieldComment,
) *ContentCommands {
	return &ContentCommands{
		CreateContent:          createContent,
		UpdateContentField:     updateContentField,
		AddContentFieldComment: addContentFieldComment,
	}
}
