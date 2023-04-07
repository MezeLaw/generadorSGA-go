package service

import "generadorSGA-go/internal/models"

type TagRepository interface {
	GetElements() *[]models.Element
}

type Service struct {
	repository TagRepository
}

func New(repository TagRepository) Service {
	return Service{repository: repository}
}

func (s *Service) GetElements() *[]models.Element {
	return s.repository.GetElements()
}
