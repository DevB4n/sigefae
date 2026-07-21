package moneda

import "sigefae/internal/db"

type Response struct {
	ID      uint   `json:"id"`
	Nombre  string `json:"nombre"`
	Codigo  string `json:"codigo"`
	Activo  bool   `json:"activo"`
}

func toResponse(moneda db.Moneda) Response {

	return Response{
		ID:      moneda.ID,
		Nombre:  moneda.Nombre,
		Codigo:  moneda.Codigo,
		Activo:  moneda.Activo,
	}
}