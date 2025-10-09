package handler

import (
	"context"
	"fmt"
	"generadorSGA-go/internal/models"
)

type TagsService interface {
	GetElements() ([]models.Element, error)
	GetElementByCode(code string) (*models.Element, error)
	SearchElements(query string) ([]models.Element, error)
}

type Handler struct {
	service TagsService
}

func New(service TagsService) Handler {
	return Handler{service: service}
}

func (h *Handler) GetTags(ctx context.Context) ([]models.Element, error) {
	elements, err := h.service.GetElements()
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	return elements, nil
}

func (h *Handler) GetTagByCode(ctx context.Context, code string) (*models.Element, error) {
	element, err := h.service.GetElementByCode(code)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag by code: %w", err)
	}

	return element, nil
}

func (h *Handler) SearchTags(ctx context.Context, query string) ([]models.Element, error) {
	elements, err := h.service.SearchElements(query)
	if err != nil {
		return nil, fmt.Errorf("failed to search tags: %w", err)
	}

	return elements, nil
}
