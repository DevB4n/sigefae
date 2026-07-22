package historial_asignacion

import (
	"time"

	"sigefae/internal/db"
)

type Response struct {
	ID                  uint       `json:"id"`
	DocumentoRadicadoID uint       `json:"documento_radicado_id"`
	UsuarioID           uint       `json:"usuario_id"`
	Desde               time.Time  `json:"desde"`
	Hasta               *time.Time `json:"hasta"`
}


func toResponse(
	historial db.HistorialAsignacion,
) Response {

	return Response{
		ID:                  historial.ID,
		DocumentoRadicadoID: historial.DocumentoRadicadoID,
		UsuarioID:           historial.UsuarioID,
		Desde:               historial.Desde,
		Hasta:               historial.Hasta,
	}
}