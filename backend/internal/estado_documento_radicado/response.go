package estado_documento_radicado

import "sigefae/internal/db"

type Response struct {
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`
	Activo bool   `json:"activo"`
}

func toResponse(estado db.EstadoDocumentoRadicado) Response {

	return Response{
		ID:     estado.ID,
		Nombre: estado.Nombre,
		Activo: estado.Activo,
	}
}
