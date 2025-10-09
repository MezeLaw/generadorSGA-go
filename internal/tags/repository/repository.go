package repository

import (
	"errors"
	"generadorSGA-go/internal/models"
)

var (
	ErrElementNotFound = errors.New("element not found")
)

type Repository struct {
	elements map[string]string
}

func New() Repository {
	elements := fillElementsDic()
	return Repository{elements: elements}
}

func (r *Repository) GetElements() []models.Element {

	elements := make([]models.Element, 0, len(r.elements))

	for code, description := range r.elements {
		element := models.Element{
			Code:        code,
			Description: description,
		}
		elements = append(elements, element)
	}

	return elements
}

func (r *Repository) GetElementByCode(code string) (*models.Element, error) {
	description, exists := r.elements[code]
	if !exists {
		return nil, ErrElementNotFound
	}

	return &models.Element{
		Code:        code,
		Description: description,
	}, nil
}

func fillElementsDic() map[string]string {

	elements := map[string]string{}
	elements["diluyenteIsoprint"] = "Diluyente Isoprint"
	elements["cleanAlPlus"] = "Clean Al Plus"
	elements["gomaArabigaNatural"] = "Goma Arabiga Natural"
	elements["limpiadorProfundoDeMantillasRevividor"] = "Limpiador Profundo de Mantillas - Revividor"
	elements["limpiadorProfundoPlusDescristalizador"] = "Limpiador Profundo Plus - Descristalizador"
	elements["reveladorPositivo"] = "Revelador Positivo"
	elements["fountGreenSF"] = "Fount Green SF"
	elements["oSSoyPlus"] = "Os Soy Plus"
	elements["universalElectroestaticEtch"] = "Universal Electroestatic Etch"

	return elements
}
