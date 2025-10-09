package handler

import (
	"context"
	"errors"
	"generadorSGA-go/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_GetTags(t *testing.T) {
	tests := []struct {
		name     string
		behavior func(m *MockTagsService)
		asserts  func(t *testing.T, elements []models.Element, err error)
	}{
		{
			name: "should return elements when service succeeds",
			behavior: func(m *MockTagsService) {
				elements := []models.Element{
					{Code: "code1", Description: "Description 1"},
					{Code: "code2", Description: "Description 2"},
				}
				m.On("GetElements").Return(elements, nil)
			},
			asserts: func(t *testing.T, elements []models.Element, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, elements)
				assert.Equal(t, 2, len(elements))
				assert.Equal(t, "code1", elements[0].Code)
				assert.Equal(t, "Description 1", elements[0].Description)
			},
		},
		{
			name: "should return error when service fails",
			behavior: func(m *MockTagsService) {
				m.On("GetElements").Return([]models.Element(nil), errors.New("service error"))
			},
			asserts: func(t *testing.T, elements []models.Element, err error) {
				assert.Error(t, err)
				assert.Nil(t, elements)
				assert.Contains(t, err.Error(), "failed to get tags")
			},
		},
		{
			name: "should return empty slice when service returns empty",
			behavior: func(m *MockTagsService) {
				m.On("GetElements").Return([]models.Element{}, nil)
			},
			asserts: func(t *testing.T, elements []models.Element, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, elements)
				assert.Empty(t, elements)
			},
		},
		{
			name: "should wrap error from service",
			behavior: func(m *MockTagsService) {
				originalErr := errors.New("no elements found")
				m.On("GetElements").Return([]models.Element(nil), originalErr)
			},
			asserts: func(t *testing.T, elements []models.Element, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to get tags")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockTagsService)
			tt.behavior(mockService)

			handler := New(mockService)
			elements, err := handler.GetTags(context.Background())

			tt.asserts(t, elements, err)
			mockService.AssertExpectations(t)
		})
	}
}

func TestHandler_GetTagByCode(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		behavior func(m *MockTagsService, code string)
		asserts  func(t *testing.T, element *models.Element, err error)
	}{
		{
			name: "should return element when code exists",
			code: "validCode",
			behavior: func(m *MockTagsService, code string) {
				element := &models.Element{
					Code:        "validCode",
					Description: "Valid Description",
				}
				m.On("GetElementByCode", code).Return(element, nil)
			},
			asserts: func(t *testing.T, element *models.Element, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, element)
				assert.Equal(t, "validCode", element.Code)
				assert.Equal(t, "Valid Description", element.Description)
			},
		},
		{
			name: "should return error when service fails",
			code: "invalidCode",
			behavior: func(m *MockTagsService, code string) {
				m.On("GetElementByCode", code).Return((*models.Element)(nil), errors.New("element not found"))
			},
			asserts: func(t *testing.T, element *models.Element, err error) {
				assert.Error(t, err)
				assert.Nil(t, element)
				assert.Contains(t, err.Error(), "failed to get tag by code")
			},
		},
		{
			name: "should handle empty code",
			code: "",
			behavior: func(m *MockTagsService, code string) {
				m.On("GetElementByCode", code).Return((*models.Element)(nil), errors.New("code cannot be empty"))
			},
			asserts: func(t *testing.T, element *models.Element, err error) {
				assert.Error(t, err)
				assert.Nil(t, element)
			},
		},
		{
			name: "should wrap service error",
			code: "someCode",
			behavior: func(m *MockTagsService, code string) {
				m.On("GetElementByCode", code).Return((*models.Element)(nil), errors.New("database error"))
			},
			asserts: func(t *testing.T, element *models.Element, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to get tag by code")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockTagsService)
			tt.behavior(mockService, tt.code)

			handler := New(mockService)
			element, err := handler.GetTagByCode(context.Background(), tt.code)

			tt.asserts(t, element, err)
			mockService.AssertExpectations(t)
		})
	}
}

func TestHandler_SearchTags(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		behavior func(m *MockTagsService, query string)
		asserts  func(t *testing.T, elements []models.Element, err error)
	}{
		{
			name:  "should return matching elements",
			query: "clean",
			behavior: func(m *MockTagsService, query string) {
				elements := []models.Element{
					{Code: "cleanAlPlus", Description: "Clean Al Plus"},
				}
				m.On("SearchElements", query).Return(elements, nil)
			},
			asserts: func(t *testing.T, elements []models.Element, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, elements)
				assert.Equal(t, 1, len(elements))
				assert.Equal(t, "cleanAlPlus", elements[0].Code)
			},
		},
		{
			name:  "should return all elements when query is empty",
			query: "",
			behavior: func(m *MockTagsService, query string) {
				elements := []models.Element{
					{Code: "code1", Description: "Description 1"},
					{Code: "code2", Description: "Description 2"},
				}
				m.On("SearchElements", query).Return(elements, nil)
			},
			asserts: func(t *testing.T, elements []models.Element, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, elements)
				assert.Equal(t, 2, len(elements))
			},
		},
		{
			name:  "should return error when no elements found",
			query: "nonexistent",
			behavior: func(m *MockTagsService, query string) {
				m.On("SearchElements", query).Return([]models.Element(nil), errors.New("no elements found"))
			},
			asserts: func(t *testing.T, elements []models.Element, err error) {
				assert.Error(t, err)
				assert.Nil(t, elements)
				assert.Contains(t, err.Error(), "failed to search tags")
			},
		},
		{
			name:  "should wrap service error",
			query: "test",
			behavior: func(m *MockTagsService, query string) {
				m.On("SearchElements", query).Return([]models.Element(nil), errors.New("search failed"))
			},
			asserts: func(t *testing.T, elements []models.Element, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to search tags")
			},
		},
		{
			name:  "should return multiple matching elements",
			query: "limpiador",
			behavior: func(m *MockTagsService, query string) {
				elements := []models.Element{
					{Code: "limpiador1", Description: "Limpiador 1"},
					{Code: "limpiador2", Description: "Limpiador 2"},
				}
				m.On("SearchElements", query).Return(elements, nil)
			},
			asserts: func(t *testing.T, elements []models.Element, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, elements)
				assert.Equal(t, 2, len(elements))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockTagsService)
			tt.behavior(mockService, tt.query)

			handler := New(mockService)
			elements, err := handler.SearchTags(context.Background(), tt.query)

			tt.asserts(t, elements, err)
			mockService.AssertExpectations(t)
		})
	}
}

func TestHandler_New(t *testing.T) {
	tests := []struct {
		name     string
		behavior func() TagsService
		asserts  func(t *testing.T, handler Handler)
	}{
		{
			name: "should create handler with service",
			behavior: func() TagsService {
				return new(MockTagsService)
			},
			asserts: func(t *testing.T, handler Handler) {
				assert.NotNil(t, handler)
				assert.NotNil(t, handler.service)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := tt.behavior()
			handler := New(service)
			tt.asserts(t, handler)
		})
	}
}

// MockTagsService is a mock implementation of TagsService interface
type MockTagsService struct {
	mock.Mock
}

func (m *MockTagsService) GetElements() ([]models.Element, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Element), args.Error(1)
}

func (m *MockTagsService) GetElementByCode(code string) (*models.Element, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Element), args.Error(1)
}

func (m *MockTagsService) SearchElements(query string) ([]models.Element, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Element), args.Error(1)
}
