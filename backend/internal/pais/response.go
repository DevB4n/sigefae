package pais

import "sigefae/internal/db"

type Response struct {
	ID      uint   `json:"id"`
	Nombre  string `json:"nombre"`
	Codigo  string `json:"codigo"`
	Activo  bool   `json:"activo"`
}

func toResponse(p db.Pais) Response {
	return Response{
		ID:      p.ID,
		Nombre:  p.Nombre,
		Codigo:  p.Codigo,
		Activo:  p.Activo,
	}
}