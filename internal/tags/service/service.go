package service

import (
	"errors"
	"generadorSGA-go/internal/models"
	"sort"
	"strings"
)

var (
	ErrNoElementsFound = errors.New("no elements found in repository")
	ErrElementNotFound = errors.New("element not found")
)

type TagRepository interface {
	GetElements() []models.Element
	GetElementByCode(code string) (*models.Element, error)
}

type Service struct {
	repository TagRepository
}

func New(repository TagRepository) Service {
	return Service{repository: repository}
}

func (s *Service) GetElements() ([]models.Element, error) {
	elements := s.repository.GetElements()

	if len(elements) == 0 {
		return nil, ErrNoElementsFound
	}

	// Sort elements by description for consistent ordering
	sort.Slice(elements, func(i, j int) bool {
		return strings.ToLower(elements[i].Description) < strings.ToLower(elements[j].Description)
	})

	return elements, nil
}

func (s *Service) GetElementByCode(code string) (*models.Element, error) {
	if code == "" {
		return nil, errors.New("code cannot be empty")
	}

	element, err := s.repository.GetElementByCode(code)
	if err != nil {
		return nil, err
	}

	return element, nil
}

func (s *Service) SearchElements(query string) ([]models.Element, error) {
	if query == "" {
		return s.GetElements()
	}

	allElements, err := s.GetElements()
	if err != nil {
		return nil, err
	}

	// Filter elements by query (case-insensitive search in code or description)
	var filtered []models.Element
	queryLower := strings.ToLower(query)

	for _, elem := range allElements {
		if strings.Contains(strings.ToLower(elem.Code), queryLower) ||
			strings.Contains(strings.ToLower(elem.Description), queryLower) {
			filtered = append(filtered, elem)
		}
	}

	if len(filtered) == 0 {
		return nil, ErrNoElementsFound
	}

	return filtered, nil
}
