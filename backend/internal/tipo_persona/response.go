package tipo_persona

import "sigefae/internal/db"

type Response struct {
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`
	Activo bool   `json:"activo"`
}

func toResponse(tipo db.TipoPersona) Response {

	return Response{
		ID:     tipo.ID,
		Nombre: tipo.Nombre,
		Activo: tipo.Activo,
	}
}
