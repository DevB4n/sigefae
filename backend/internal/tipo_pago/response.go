package tipo_pago

import "sigefae/internal/db"

type Response struct {
	ID      uint   `json:"id"`
	Nombre  string `json:"nombre"`
	Activo  bool   `json:"activo"`
}

func toResponse(tipo db.TipoPago) Response {

	return Response{
		ID:      tipo.ID,
		Nombre:  tipo.Nombre,
		Activo:  tipo.Activo,
	}
}