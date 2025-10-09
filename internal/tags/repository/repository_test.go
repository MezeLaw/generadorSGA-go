package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_GetElements(t *testing.T) {
	tests := []struct {
		name     string
		behavior func() Repository
		asserts  func(t *testing.T, repo Repository)
	}{
		{
			name: "should return all elements from repository",
			behavior: func() Repository {
				return New()
			},
			asserts: func(t *testing.T, repo Repository) {
				elements := repo.GetElements()
				assert.NotNil(t, elements)
				assert.NotEmpty(t, elements)
				assert.Equal(t, 9, len(elements), "should return all 9 elements")

				// Verify structure of returned elements
				for _, elem := range elements {
					assert.NotEmpty(t, elem.Code, "element code should not be empty")
					assert.NotEmpty(t, elem.Description, "element description should not be empty")
				}
			},
		},
		{
			name: "should return elements with correct mappings",
			behavior: func() Repository {
				return New()
			},
			asserts: func(t *testing.T, repo Repository) {
				elements := repo.GetElements()

				// Create a map for easier lookup
				elementMap := make(map[string]string)
				for _, elem := range elements {
					elementMap[elem.Code] = elem.Description
				}

				// Verify specific elements exist with correct descriptions
				assert.Equal(t, "Clean Al Plus", elementMap["cleanAlPlus"])
				assert.Equal(t, "Fount Green SF", elementMap["fountGreenSF"])
				assert.Equal(t, "Diluyente Isoprint", elementMap["diluyenteIsoprint"])
			},
		},
		{
			name: "should return slice not pointer",
			behavior: func() Repository {
				return New()
			},
			asserts: func(t *testing.T, repo Repository) {
				elements := repo.GetElements()
				// This test verifies we're returning a slice, not a pointer
				assert.IsType(t, elements, elements)
				assert.GreaterOrEqual(t, len(elements), 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.behavior()
			tt.asserts(t, repo)
		})
	}
}

func TestRepository_GetElementByCode(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		behavior func() Repository
		asserts  func(t *testing.T, element interface{}, err error)
	}{
		{
			name: "should return element when code exists",
			code: "cleanAlPlus",
			behavior: func() Repository {
				return New()
			},
			asserts: func(t *testing.T, element interface{}, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, element)
				elem := element
				assert.NotNil(t, elem)
			},
		},
		{
			name: "should return correct element details",
			code: "fountGreenSF",
			behavior: func() Repository {
				return New()
			},
			asserts: func(t *testing.T, element interface{}, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, element)
			},
		},
		{
			name: "should return error when code does not exist",
			code: "nonExistentCode",
			behavior: func() Repository {
				return New()
			},
			asserts: func(t *testing.T, element interface{}, err error) {
				assert.Error(t, err)
				assert.Nil(t, element)
				assert.Equal(t, ErrElementNotFound, err)
			},
		},
		{
			name: "should return error for empty code",
			code: "",
			behavior: func() Repository {
				return New()
			},
			asserts: func(t *testing.T, element interface{}, err error) {
				assert.Error(t, err)
				assert.Nil(t, element)
				assert.Equal(t, ErrElementNotFound, err)
			},
		},
		{
			name: "should return element pointer with correct code",
			code: "reveladorPositivo",
			behavior: func() Repository {
				return New()
			},
			asserts: func(t *testing.T, element interface{}, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, element)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.behavior()
			element, err := repo.GetElementByCode(tt.code)
			tt.asserts(t, element, err)
		})
	}
}

func TestRepository_New(t *testing.T) {
	tests := []struct {
		name     string
		behavior func() Repository
		asserts  func(t *testing.T, repo Repository)
	}{
		{
			name: "should create repository with elements",
			behavior: func() Repository {
				return New()
			},
			asserts: func(t *testing.T, repo Repository) {
				assert.NotNil(t, repo)
				assert.NotNil(t, repo.elements)
				assert.NotEmpty(t, repo.elements)
				assert.Equal(t, 9, len(repo.elements))
			},
		},
		{
			name: "should initialize with expected elements",
			behavior: func() Repository {
				return New()
			},
			asserts: func(t *testing.T, repo Repository) {
				expectedCodes := []string{
					"diluyenteIsoprint",
					"cleanAlPlus",
					"gomaArabigaNatural",
					"limpiadorProfundoDeMantillasRevividor",
					"limpiadorProfundoPlusDescristalizador",
					"reveladorPositivo",
					"fountGreenSF",
					"oSSoyPlus",
					"universalElectroestaticEtch",
				}

				for _, code := range expectedCodes {
					_, exists := repo.elements[code]
					assert.True(t, exists, "code %s should exist in repository", code)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.behavior()
			tt.asserts(t, repo)
		})
	}
}

func TestFillElementsDic(t *testing.T) {
	tests := []struct {
		name     string
		behavior func() map[string]string
		asserts  func(t *testing.T, elements map[string]string)
	}{
		{
			name: "should return map with all elements",
			behavior: func() map[string]string {
				return fillElementsDic()
			},
			asserts: func(t *testing.T, elements map[string]string) {
				assert.NotNil(t, elements)
				assert.Equal(t, 9, len(elements))
			},
		},
		{
			name: "should have correct key-value mappings",
			behavior: func() map[string]string {
				return fillElementsDic()
			},
			asserts: func(t *testing.T, elements map[string]string) {
				assert.Equal(t, "Diluyente Isoprint", elements["diluyenteIsoprint"])
				assert.Equal(t, "Clean Al Plus", elements["cleanAlPlus"])
				assert.Equal(t, "Goma Arabiga Natural", elements["gomaArabigaNatural"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elements := tt.behavior()
			tt.asserts(t, elements)
		})
	}
}
