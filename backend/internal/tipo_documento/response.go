package tipo_documento

import "sigefae/internal/db"

type Response struct{
	ID uint `json:"id"`
	Nombre string `json:"nombre"`
	Activo bool `json:"activo"`
}

func toResponse(tipo db.TipoDocumento) Response {

	return Response{
		ID: tipo.ID,
		Nombre: tipo.Nombre,
		Activo: tipo.Activo,
	}
}