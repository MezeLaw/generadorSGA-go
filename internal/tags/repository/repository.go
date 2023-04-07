package repository

import "generadorSGA-go/internal/models"

type Repository struct {
	elements map[string]string
}

func New() Repository {
	elements := fillElementsDic()
	return Repository{elements: elements}
}

func (r *Repository) GetElements() *[]models.Element {

	elements := []models.Element{}

	for code, description := range r.elements {
		element := models.Element{
			Code:        code,
			Description: description,
		}
		elements = append(elements, element)
	}

	return &elements
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
