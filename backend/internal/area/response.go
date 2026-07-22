package area

import "sigefae/internal/db"

type Response struct {
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`
	Activo bool   `json:"activo"`
}

func toResponse(area db.Area) Response {

	return Response{
		ID:     area.ID,
		Nombre: area.Nombre,
		Activo: area.Activo,
	}
}
