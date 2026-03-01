package dto

import "ownned/internal/domain"

type FileNodeDTO struct {
	domain.Node
	Doc domain.Doc
}

type FolderNodeDTO struct {
	domain.Node
	Children []domain.Node
}
