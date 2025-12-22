package model

import "ownned/internal/domain"

type FileNodeDTO struct {
	*domain.Node
	Docs []domain.Doc
}

type FolderNodeDTO struct {
	*domain.Node
	Children []domain.Node
}
