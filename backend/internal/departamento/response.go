package departamento

import "sigefae/internal/db"

type PaisResponse struct {
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`
	Codigo string `json:"codigo"`
}

type Response struct {
	ID     uint          `json:"id"`
	Nombre string        `json:"nombre"`

	PaisID uint          `json:"pais_id"`
	Pais   *PaisResponse `json:"pais,omitempty"`

	Activo bool `json:"activo"`
}

func toResponse(departamento db.Departamento) Response {

	response := Response{
		ID:     departamento.ID,
		Nombre: departamento.Nombre,
		PaisID: departamento.PaisID,
		Activo: departamento.Activo,
	}

	if departamento.Pais != nil {
		response.Pais = &PaisResponse{
			ID:     departamento.Pais.ID,
			Nombre: departamento.Pais.Nombre,
			Codigo: departamento.Pais.Codigo,
		}
	}

	return response
}