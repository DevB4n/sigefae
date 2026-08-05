package norma_reparto

import "sigefae/internal/db"

type Response struct {
	ID           uint    `json:"id"`
	Codigo       string  `json:"codigo"`
	Nombre       string  `json:"nombre"`
	Sucursal     string  `json:"sucursal"`
	Departamento string  `json:"departamento"`
	Tipo         *string `json:"tipo,omitempty"`
	TarifaIva    *string `json:"tarifa_iva,omitempty"`
	Activo       bool    `json:"activo"`
}

func toResponse(n db.NormaReparto) Response {
	return Response{
		ID:           n.ID,
		Codigo:       n.Codigo,
		Nombre:       n.Nombre,
		Sucursal:     n.Sucursal,
		Departamento: n.Departamento,
		Tipo:         n.Tipo,
		TarifaIva:    n.TarifaIva,
		Activo:       n.Activo,
	}
}