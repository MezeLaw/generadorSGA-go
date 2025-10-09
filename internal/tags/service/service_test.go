package service

import (
	"generadorSGA-go/internal/models"
	"testing"
)

// Mock repository for testing
type mockRepository struct {
	elements map[string]string
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		elements: map[string]string{
			"cleanAlPlus":     "Clean Al Plus",
			"fountGreenSF":    "Fount Green SF",
			"reveladorPositivo": "Revelador Positivo",
		},
	}
}

func (m *mockRepository) GetElements() []models.Element {
	elements := make([]models.Element, 0, len(m.elements))
	for code, description := range m.elements {
		elements = append(elements, models.Element{
			Code:        code,
			Description: description,
		})
	}
	return elements
}

func (m *mockRepository) GetElementByCode(code string) (*models.Element, error) {
	description, exists := m.elements[code]
	if !exists {
		return nil, ErrElementNotFound
	}
	return &models.Element{
		Code:        code,
		Description: description,
	}, nil
}

func TestGetElements(t *testing.T) {
	repo := newMockRepository()
	service := New(repo)

	elements, err := service.GetElements()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(elements) != 3 {
		t.Errorf("Expected 3 elements, got %d", len(elements))
	}

	// Verify elements are sorted by description
	if len(elements) > 1 {
		for i := 0; i < len(elements)-1; i++ {
			if elements[i].Description > elements[i+1].Description {
				t.Errorf("Elements not sorted: %s comes after %s", elements[i].Description, elements[i+1].Description)
			}
		}
	}
}

func TestGetElementByCode(t *testing.T) {
	repo := newMockRepository()
	service := New(repo)

	tests := []struct {
		name        string
		code        string
		expectError bool
		expectDesc  string
	}{
		{
			name:        "Valid code",
			code:        "cleanAlPlus",
			expectError: false,
			expectDesc:  "Clean Al Plus",
		},
		{
			name:        "Invalid code",
			code:        "nonexistent",
			expectError: true,
		},
		{
			name:        "Empty code",
			code:        "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			element, err := service.GetElementByCode(tt.code)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if element.Description != tt.expectDesc {
					t.Errorf("Expected description %s, got %s", tt.expectDesc, element.Description)
				}
			}
		})
	}
}

func TestSearchElements(t *testing.T) {
	repo := newMockRepository()
	service := New(repo)

	tests := []struct {
		name          string
		query         string
		expectError   bool
		expectedCount int
	}{
		{
			name:          "Search by partial description",
			query:         "Clean",
			expectError:   false,
			expectedCount: 1,
		},
		{
			name:          "Search by code",
			query:         "fount",
			expectError:   false,
			expectedCount: 1,
		},
		{
			name:          "No results",
			query:         "nonexistent",
			expectError:   true,
			expectedCount: 0,
		},
		{
			name:          "Empty query returns all",
			query:         "",
			expectError:   false,
			expectedCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elements, err := service.SearchElements(tt.query)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if len(elements) != tt.expectedCount {
					t.Errorf("Expected %d elements, got %d", tt.expectedCount, len(elements))
				}
			}
		})
	}
}

func TestGetElementsEmptyRepository(t *testing.T) {
	emptyRepo := &mockRepository{elements: map[string]string{}}
	service := New(emptyRepo)

	_, err := service.GetElements()
	if err != ErrNoElementsFound {
		t.Errorf("Expected ErrNoElementsFound, got %v", err)
	}
}
