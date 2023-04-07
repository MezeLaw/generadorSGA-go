package handler

import (
	"context"
	"errors"
	"generadorSGA-go/internal/models"
)

type TagsService interface {
	GetElements() *[]models.Element
}

type Handler struct {
	service TagsService
}

func New(service TagsService) Handler {
	return Handler{service: service}
}

func (h *Handler) GetTags(ctx context.Context) (*[]models.Element, error) {
	elements := h.service.GetElements()

	if len(*elements) < 1 {
		return nil, errors.New("no tags found")
	}

	return elements, nil
}

func (h *Handler) getTagByCode(ctx context.Context, code string) (*models.Element, error) {
	return nil, nil
}
