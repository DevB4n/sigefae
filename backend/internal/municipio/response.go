package municipio

import "sigefae/internal/db"

type DepartamentoResponse struct {
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`
}

type Response struct {
	ID             uint                  `json:"id"`
	Nombre         string                `json:"nombre"`

	DepartamentoID uint                  `json:"departamento_id"`
	Departamento   *DepartamentoResponse `json:"departamento,omitempty"`

	Activo         bool                  `json:"activo"`
}

func toResponse(municipio db.Municipio) Response {

	response := Response{
		ID:             municipio.ID,
		Nombre:         municipio.Nombre,
		DepartamentoID: municipio.DepartamentoID,
		Activo:         municipio.Activo,
	}

	if municipio.Departamento != nil {
		response.Departamento = &DepartamentoResponse{
			ID:     municipio.Departamento.ID,
			Nombre: municipio.Departamento.Nombre,
		}
	}

	return response
}