package actividad_economica

import "sigefae/internal/db"

type Response struct {
	ID      uint   `json:"id"`
	Nombre  string `json:"nombre"`
	Codigo  string `json:"codigo"`
	Activo  bool   `json:"activo"`
}

func toResponse(a db.ActividadEconomica) Response {

	return Response{
		ID:      a.ID,
		Nombre:  a.Nombre,
		Codigo:  a.Codigo,
		Activo:  a.Activo,
	}
}