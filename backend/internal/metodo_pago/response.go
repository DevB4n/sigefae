package metodo_pago

import "sigefae/internal/db"

type Response struct {
	ID         uint   `json:"id"`
	TipoPagoID uint   `json:"tipo_pago_id"`
	TipoPago   string `json:"tipo_pago"`
	Nombre     string `json:"nombre"`
	Activo     bool   `json:"activo"`
}

func toResponse(m db.MetodoPago) Response {

	response := Response{
		ID:         m.ID,
		TipoPagoID: m.TipoPagoID,
		Nombre:     m.Nombre,
		Activo:     m.Activo,
	}

	if m.TipoPago != nil {
		response.TipoPago = m.TipoPago.Nombre
	}

	return response
}
