package service

import (
	"errors"
	"generadorSGA-go/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_GetElements(t *testing.T) {
	tests := []struct {
		name     string
		behavior func(m *MockTagRepository)
		asserts  func(t *testing.T, elements []models.Element, err error)
	}{
		{
			name: "should return sorted elements from repository",
			behavior: func(m *MockTagRepository) {
				elements := []models.Element{
					{Code: "code3", Description: "Zebra Description"},
					{Code: "code1", Description: "Alpha Description"},
					{Code: "code2", Description: "Beta Description"},
				}
				m.On("GetElements").Return(elements)
			},
			asserts: func(t *testing.T, elements []models.Element, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, elements)
				assert.Equal(t, 3, len(elements))
				// Verify sorting by description
				assert.Equal(t, "Alpha Description", elements[0].Description)
				assert.Equal(t, "Beta Description", elements[1].Description)
				assert.Equal(t, "Zebra Description", elements[2].Description)
			},
		},
		{
			name: "should return error when repository is empty",
			behavior: func(m *MockTagRepository) {
				m.On("GetElements").Return([]models.Element{})
			},
			asserts: func(t *testing.T, elements []models.Element, err error) {
				assert.Error(t, err)
				assert.Equal(t, ErrNoElementsFound, err)
				assert.Nil(t, elements)
			},
		},
		{
			name: "should sort case-insensitively",
			behavior: func(m *MockTagRepository) {
				elements := []models.Element{
					{Code: "code1", Description: "apple"},
					{Code: "code2", Description: "Banana"},
					{Code: "code3", Description: "cherry"},
				}
				m.On("GetElements").Return(elements)
			},
			asserts: func(t *testing.T, elements []models.Element, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 3, len(elements))
				assert.Equal(t, "apple", elements[0].Description)
				assert.Equal(t, "Banana", elements[1].Description)
				assert.Equal(t, "cherry", elements[2].Description)
			},
		},
		{
			name: "should handle single element",
			behavior: func(m *MockTagRepository) {
				elements := []models.Element{
					{Code: "code1", Description: "Single Description"},
				}
				m.On("GetElements").Return(elements)
			},
			asserts: func(t *testing.T, elements []models.Element, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 1, len(elements))
				assert.Equal(t, "Single Description", elements[0].Description)
			},
		},
		{
			name: "should maintain element integrity while sorting",
			behavior: func(m *MockTagRepository) {
				elements := []models.Element{
					{Code: "code2", Description: "Zebra"},
					{Code: "code1", Description: "Alpha"},
				}
				m.On("GetElements").Return(elements)
			},
			asserts: func(t *testing.T, elements []models.Element, err error) {
				assert.NoError(t, err)
				// Verify code-description pairing is maintained
				for _, elem := range elements {
					if elem.Description == "Alpha" {
						assert.Equal(t, "code1", elem.Code)
					}
					if elem.Description == "Zebra" {
						assert.Equal(t, "code2", elem.Code)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTagRepository)
			tt.behavior(mockRepo)

			service := New(mockRepo)
			elements, err := service.GetElements()

			tt.asserts(t, elements, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetElementByCode(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		behavior func(m *MockTagRepository, code string)
		asserts  func(t *testing.T, element *models.Element, err error)
	}{
		{
			name: "should return element when code exists",
			code: "validCode",
			behavior: func(m *MockTagRepository, code string) {
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
			name: "should return error when code is empty",
			code: "",
			behavior: func(m *MockTagRepository, code string) {
				// Repository should not be called for empty code
			},
			asserts: func(t *testing.T, element *models.Element, err error) {
				assert.Error(t, err)
				assert.Nil(t, element)
				assert.Contains(t, err.Error(), "code cannot be empty")
			},
		},
		{
			name: "should return error when code does not exist",
			code: "nonExistentCode",
			behavior: func(m *MockTagRepository, code string) {
				m.On("GetElementByCode", code).Return((*models.Element)(nil), ErrElementNotFound)
			},
			asserts: func(t *testing.T, element *models.Element, err error) {
				assert.Error(t, err)
				assert.Nil(t, element)
				assert.Equal(t, ErrElementNotFound, err)
			},
		},
		{
			name: "should propagate repository errors",
			code: "errorCode",
			behavior: func(m *MockTagRepository, code string) {
				m.On("GetElementByCode", code).Return((*models.Element)(nil), errors.New("repository error"))
			},
			asserts: func(t *testing.T, element *models.Element, err error) {
				assert.Error(t, err)
				assert.Nil(t, element)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTagRepository)
			tt.behavior(mockRepo, tt.code)

			service := New(mockRepo)
			element, err := service.GetElementByCode(tt.code)

			tt.asserts(t, element, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_SearchElements(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		behavior func(m *MockTagRepository)
		asserts  func(t *testing.T, elements []models.Element, err error)
	}{
		{
			name:  "should return all elements when query is empty",
			query: "",
			behavior: func(m *MockTagRepository) {
				elements := []models.Element{
					{Code: "code1", Description: "Description 1"},
					{Code: "code2", Description: "Description 2"},
				}
				m.On("GetElements").Return(elements)
			},
			asserts: func(t *testing.T, elements []models.Element, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 2, len(elements))
			},
		},
		{
			name:  "should filter by code case-insensitive",
			query: "CLEAN",
			behavior: func(m *MockTagRepository) {
				elements := []models.Element{
					{Code: "cleanAlPlus", Description: "Clean Al Plus"},
					{Code: "fountGreenSF", Description: "Fount Green SF"},
				}
				m.On("GetElements").Return(elements)
			},
			asserts: func(t *testing.T, elements []models.Element, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 1, len(elements))
				assert.Equal(t, "cleanAlPlus", elements[0].Code)
			},
		},
		{
			name:  "should filter by description case-insensitive",
			query: "green",
			behavior: func(m *MockTagRepository) {
				elements := []models.Element{
					{Code: "cleanAlPlus", Description: "Clean Al Plus"},
					{Code: "fountGreenSF", Description: "Fount Green SF"},
				}
				m.On("GetElements").Return(elements)
			},
			asserts: func(t *testing.T, elements []models.Element, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 1, len(elements))
				assert.Equal(t, "Fount Green SF", elements[0].Description)
			},
		},
		{
			name:  "should return error when no matches found",
			query: "nonexistent",
			behavior: func(m *MockTagRepository) {
				elements := []models.Element{
					{Code: "code1", Description: "Description 1"},
				}
				m.On("GetElements").Return(elements)
			},
			asserts: func(t *testing.T, elements []models.Element, err error) {
				assert.Error(t, err)
				assert.Equal(t, ErrNoElementsFound, err)
				assert.Nil(t, elements)
			},
		},
		{
			name:  "should return multiple matches",
			query: "limpiador",
			behavior: func(m *MockTagRepository) {
				elements := []models.Element{
					{Code: "limpiador1", Description: "Limpiador Profundo"},
					{Code: "limpiador2", Description: "Limpiador Plus"},
					{Code: "other", Description: "Other Product"},
				}
				m.On("GetElements").Return(elements)
			},
			asserts: func(t *testing.T, elements []models.Element, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 2, len(elements))
			},
		},
		{
			name:  "should handle partial matches",
			query: "Al",
			behavior: func(m *MockTagRepository) {
				elements := []models.Element{
					{Code: "cleanAlPlus", Description: "Clean Al Plus"},
					{Code: "fountGreenSF", Description: "Fount Green SF"},
				}
				m.On("GetElements").Return(elements)
			},
			asserts: func(t *testing.T, elements []models.Element, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 1, len(elements))
				assert.Contains(t, elements[0].Code, "Al")
			},
		},
		{
			name:  "should propagate error from GetElements",
			query: "test",
			behavior: func(m *MockTagRepository) {
				m.On("GetElements").Return([]models.Element{})
			},
			asserts: func(t *testing.T, elements []models.Element, err error) {
				assert.Error(t, err)
				assert.Equal(t, ErrNoElementsFound, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTagRepository)
			tt.behavior(mockRepo)

			service := New(mockRepo)
			elements, err := service.SearchElements(tt.query)

			tt.asserts(t, elements, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_New(t *testing.T) {
	tests := []struct {
		name     string
		behavior func() TagRepository
		asserts  func(t *testing.T, service Service)
	}{
		{
			name: "should create service with repository",
			behavior: func() TagRepository {
				return new(MockTagRepository)
			},
			asserts: func(t *testing.T, service Service) {
				assert.NotNil(t, service)
				assert.NotNil(t, service.repository)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.behavior()
			service := New(repo)
			tt.asserts(t, service)
		})
	}
}

// MockTagRepository is a mock implementation of TagRepository interface
type MockTagRepository struct {
	mock.Mock
}

func (m *MockTagRepository) GetElements() []models.Element {
	args := m.Called()
	return args.Get(0).([]models.Element)
}

func (m *MockTagRepository) GetElementByCode(code string) (*models.Element, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Element), args.Error(1)
}
